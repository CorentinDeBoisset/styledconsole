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

// FormatAndWrap find all tags and replace them with the correct escape sequences,
// and adds newlines when necessary to ensure the output is fine in a given terminal
func FormatAndWrap(text string, width int) string {
	var offset int
	var output string
	currentLineLength := new(int)
	*currentLineLength = 0

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
			formatStringWithStyle(text[offset:tagIndexes[0]], &output, width, currentLineLength, styleStack),
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
					formatStringWithStyle(text[tagIndexes[0]:tagIndexes[1]], &output, width, currentLineLength, styleStack),
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
		formatStringWithStyle(text[offset:], &output, width, currentLineLength, styleStack),
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
func formatStringWithStyle(text string, current *string, width int, currentLineLength *int, stack style.OutputStyleStack) string {
	if text == `` {
		return ``
	}
	if width == 0 || currentLineLength == nil {
		return stack.GetCurrent().Apply(text)
	}

	if *currentLineLength == 0 {
		text = strings.TrimLeft(text, " ")
	}

	// The prefix is aimed at filling the current line with just enought characters
	prefix := ``
	if *currentLineLength > 0 {
		prefix = fmt.Sprintf("%s\n", strings.TrimRight(getSubstring(text, 0, width-*currentLineLength), " "))
		text = getSubstring(text, width-*currentLineLength, len(text))
	}

	// Then we cut the text into width-sized elements
	lineSplitRegexp := regexp.MustCompile(fmt.Sprintf(`([^\n]{%d}) *`, width))
	textHasNewLine := len(text) > 0 && text[len(text)-1:] == "\n"
	text = fmt.Sprintf("%s%s", prefix, lineSplitRegexp.ReplaceAllString(text, "$1\n"))

	// Merge all line endings at the end of the text together
	if textHasNewLine {
		text = fmt.Sprintf("%s\n", strings.TrimRight(text, "\n"))
	}

	// If there is no started line, and the last item in the current string is not a \n, we add one before the new text
	if *currentLineLength == 0 && len(*current) > 0 && (*current)[len(*current)-1:] != "\n" {
		text = fmt.Sprintf("\n%s", text)
	}

	lines := strings.Split(text, "\n")
	for i, line := range lines {
		// First we set currentLineLength
		*currentLineLength = *currentLineLength + len(line)
		if width <= *currentLineLength {
			*currentLineLength = 0
		}

		// Then we decorate each line
		lines[i] = stack.GetCurrent().Apply(line)
	}

	return strings.Join(lines, "\n")
}
