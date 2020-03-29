package internal

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/coreoas/styledconsole/internal/style"
)

var (
	escapeRegexp    = regexp.MustCompile(`[^\\]?<`)
	tagRegexp       = regexp.MustCompile(`(?i)<([a-z][^<>]*|/([a-z][^<>]*)?)>`)
	styleRegexp     = regexp.MustCompile(`([^=]+)=([^;]+)(;|$)`)
	separatorRegexp = regexp.MustCompile(`([^,;]+)`)
	lineEndRegexp   = regexp.MustCompile(` *(\r?\n)`)

	// Some useful pre-defined styles
	errorStyle    style.OutputStyle
	infoStyle     style.OutputStyle
	commentStyle  style.OutputStyle
	questionStyle style.OutputStyle
)

// Escape will prepend all '<' with a backslash
func Escape(text string) string {
	return escapeRegexp.ReplaceAllString(text, `$1\<`)
}

// EscapeTrailingBackslash removes any trailing backslashes while keeping the length of the string
func EscapeTrailingBackslash(text string) string {
	newText := text
	textLen := len(newText)
	if textLen > 0 && newText[textLen-1:] == `\` {
		newText = strings.TrimRight(text, `\`)
		newText = strings.ReplaceAll(newText, "\x00", ``)
		newText = fmt.Sprintf("%s%s", newText, strings.Repeat("\x00", textLen-len(newText)))
	}

	return newText
}

// FormatText find all tags and replace them with the correct escape sequences,
// and adds newlines when necessary to ensure the output is fine in a given terminal
func FormatText(text string, width int) string {
	var offset int
	var output string
	currentLineLength := 0

	tagMatches := tagRegexp.FindAllSubmatchIndex([]byte(text), -1)
	styleStack := style.OutputStyleStack{}

	for _, tagIndexes := range tagMatches {
		if tagIndexes[0] == 0 && text[len(text)-1:] == `\` {
			continue
		}

		// Write text up to the tag
		output = fmt.Sprintf(
			"%s%s",
			output,
			formatStringWithStyle(text[offset:tagIndexes[0]], width, &currentLineLength, styleStack),
		)
		offset = tagIndexes[1]

		// Opening tag ?
		tagName := ``
		openingTag := text[tagIndexes[2]:tagIndexes[3]] != `/`
		if openingTag {
			tagName = text[tagIndexes[2]:tagIndexes[3]]
		} else if len(tagIndexes) >= 6 && tagIndexes[4] > 0 && tagIndexes[5] > 0 {
			tagName = text[tagIndexes[4]:tagIndexes[5]]
		}

		if !openingTag && len(tagName) == 0 {
			// tag is </>
			styleStack.PopCurrent()
		} else {
			style := extractStyle(tagName)
			if style == nil {
				// We detected a tag incorrectly, we write the text in the regex
				output = fmt.Sprintf(
					"%s%s",
					output,
					formatStringWithStyle(text[tagIndexes[0]:tagIndexes[1]], width, &currentLineLength, styleStack),
				)
			} else if openingTag {
				styleStack.Push(*style)
			} else {
				styleStack.Pop(*style)
			}
		}
	}

	// Write the end of the text
	output = fmt.Sprintf(
		"%s%s",
		output,
		formatStringWithStyle(text[offset:], width, &currentLineLength, styleStack),
	)

	output = strings.ReplaceAll(output, "\x00", `\`)
	output = strings.ReplaceAll(output, `\<`, `<`)

	return output
}

func getSubstring(s string, start int, end int) string {
	if start > len(s) {
		return ``
	}
	if end > len(s) {
		return s[start:len(s)]
	}

	return s[start:end]
}

func extractStyle(tagName string) *style.OutputStyle {
	styleMatches := styleRegexp.FindAllSubmatchIndex([]byte(tagName), -1)
	if len(styleMatches) == 0 {
		return nil
	}

	extractedStyle := new(style.OutputStyle)
	for _, attrMatches := range styleMatches {
		styleName := strings.ToLower(tagName[attrMatches[2]:attrMatches[3]])
		styleValue := tagName[attrMatches[4]:attrMatches[5]]

		if `fg` == styleName {
			extractedStyle.Foreground = strings.ToLower(styleValue)
		} else if `bg` == styleName {
			extractedStyle.Background = strings.ToLower(styleValue)
		} else if `href` == styleName {
			extractedStyle.Href = styleValue
		} else if `options` == styleName {
			separatorMatches := separatorRegexp.FindAllSubmatchIndex([]byte(strings.ToLower(styleValue)), -1)
			for _, separatorIndexes := range separatorMatches {
				extractedStyle.Options = append(extractedStyle.Options, styleValue[separatorIndexes[2]:separatorIndexes[3]])
			}
		} else {
			// If there is an unknown attribute, the whole style is voided
			return nil
		}
	}

	return extractedStyle
}

// This function is pretty bad, it should be much more clean and thoroughly tested
func formatStringWithStyle(text string, width int, currentLineLength *int, stack style.OutputStyleStack) string {
	// First, handle invalid argument cases
	if text == `` {
		return ``
	}
	if width == 0 || currentLineLength == nil {
		return stack.GetCurrent().Apply(text)
	}

	// First cleanup text and replace line endings with \n, then split the lines
	sourceLines := strings.Split(lineEndRegexp.ReplaceAllString(text, "\n"), "\n")

	var splitLines []string

	if *currentLineLength < 0 {
		*currentLineLength = 0
	} else if *currentLineLength > width {
		splitLines = append(splitLines, "")
		*currentLineLength = width
	} else if *currentLineLength > 0 && *currentLineLength+len(sourceLines[0]) > width {
		// If required, split the first line in two
		splitLines = append(
			splitLines,
			getSubstring(sourceLines[0], 0, width-*currentLineLength),
		)
		sourceLines[0] = getSubstring(sourceLines[0], width-*currentLineLength, len(sourceLines[0]))
	}

	// Then split all the other lines.
	for _, line := range sourceLines {
		nbSublines := len(line) / width
		for j := 0; j <= nbSublines; j++ {
			splitLines = append(splitLines, getSubstring(line, j*width, (j+1)*width))
		}
	}

	// Remove all empty elements (=newLines) but one at the end of splitLines
	textHasNewLine := false
	for i := len(splitLines) - 1; i >= 0; i-- {
		if len(splitLines[i]) > 0 {
			splitLines = splitLines[0 : i+1]
			break
		}
		textHasNewLine = true
	}
	if textHasNewLine {
		splitLines = append(splitLines, "")
	}

	// Fill the lines with spaces
	for i, line := range splitLines[:len(splitLines)-1] {
		// Special case for the first line that has to takes into account *currentLineLength
		if i == 0 && (len(line)+*currentLineLength) < width {
			splitLines[i] = line + strings.Repeat(" ", width-len(line))
		} else if i > 0 && len(line) < width {
			splitLines[i] = line + strings.Repeat(" ", width-len(line))
		}
	}

	for i, line := range splitLines {
		// First we set currentLineLength
		*currentLineLength = *currentLineLength + len(line)
		if *currentLineLength >= width {
			*currentLineLength = 0
		}

		if len(line) > 0 {
			// Then we decorate each line
			splitLines[i] = stack.GetCurrent().Apply(line)
		}
	}

	return strings.Join(splitLines, "\n")
}
