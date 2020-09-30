package internal

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/corentindeboisset/styledconsole/internal/style"
)

var (
	escapeRegexp  = regexp.MustCompile(`[^\\]?<`)
	tagRegexp     = regexp.MustCompile(`(?i)<([a-z][^<>]*|/([a-z][^<>]*)?)>`)
	lineEndRegexp = regexp.MustCompile(` *(\r?\n)`)
)

// Escape will prepend all '<' with a backslash
func Escape(text string) string {
	return escapeRegexp.ReplaceAllString(text, `$1\<`)
}

// EscapeTrailingBackslash removes any trailing backslashes while keeping the length of the string
func EscapeTrailingBackslash(text string) string {
	textLen := len(text)
	if textLen > 0 && text[textLen-1:] == `\` {
		newText := strings.TrimRight(text, `\`)
		newText = strings.TrimRight(newText, "\x00")
		newText = fmt.Sprintf("%s%s", newText, strings.Repeat("\x00", textLen-len(newText)))

		return newText
	}

	return text
}

// FormatText find all tags and replace them with the correct escape sequences,
// and adds newlines when necessary to ensure the output is fine in a given terminal
func FormatText(text string, width int, baseStyle *style.OutputStyle) []string {
	var offset int
	var styleStack style.OutputStyleStack

	output := []string{""}
	currentLineLength := 0

	tagMatches := tagRegexp.FindAllSubmatchIndex([]byte(text), -1)
	if baseStyle != nil {
		styleStack = style.OutputStyleStack{BaseStyle: baseStyle}
	} else {
		styleStack = style.OutputStyleStack{}
	}

	for _, tagIndexes := range tagMatches {
		if tagIndexes[0] == 0 && text[len(text)-1] == '\\' {
			continue
		}

		// Write text up to the tag
		addStringWithStyle(text[offset:tagIndexes[0]], width, &output, &currentLineLength, styleStack)
		offset = tagIndexes[1]

		// Opening tag ?
		tagName := ``
		// openingTag := tagIndexes
		openingTag := text[tagIndexes[2]:tagIndexes[3]][0] != '/'
		if openingTag {
			tagName = text[tagIndexes[2]:tagIndexes[3]]
		} else if tagIndexes[4] > 0 && tagIndexes[5] > 0 {
			tagName = text[tagIndexes[4]:tagIndexes[5]]
		}

		if !openingTag && len(tagName) == 0 {
			// tag is </>
			styleStack.PopCurrent()
		} else {
			style := style.NewOutputStyle(tagName)
			if style == nil {
				// We detected a tag incorrectly, we write the text in the regex
				addStringWithStyle(text[tagIndexes[0]:tagIndexes[1]], width, &output, &currentLineLength, styleStack)
			} else if openingTag {
				styleStack.Push(*style)
			} else {
				styleStack.Pop(*style)
			}
		}
	}

	// Write the end of the text
	addStringWithStyle(text[offset:], width, &output, &currentLineLength, styleStack)

	for i, line := range output {
		output[i] = strings.ReplaceAll(line, "\x00", `\`)
		output[i] = strings.ReplaceAll(line, `\<`, `<`)
	}

	return output
}

func getSubstring(s string, start int, end int) string {
	if start > len(s) {
		return ``
	}
	if end > len(s) {
		return s[start:]
	}

	return s[start:end]
}

// This function is pretty bad, it should be much more clean and thoroughly tested
func addStringWithStyle(text string, width int, output *[]string, lastLineLength *int, stack style.OutputStyleStack) {
	// First, handle invalid argument cases
	if text == `` || width == 0 || output == nil {
		return
	}

	// First cleanup text and replace line endings with \n, then split the lines
	sourceLines := strings.Split(lineEndRegexp.ReplaceAllString(text, "\n"), "\n")

	var splitLines []string

	if *lastLineLength < 0 {
		// The lastLineLength was invalid, we try to infer it from *output
		// It's not optimal because the elements of output contain escape codes
		*lastLineLength = 0
		if len(*output) > 0 {
			*lastLineLength = len((*output)[len(*output)-1])
		}
	} else if *lastLineLength > width {
		splitLines = append(splitLines, "")
		*lastLineLength = width
	} else if *lastLineLength > 0 && *lastLineLength+len(sourceLines[0]) > width {
		// If required, split the first line in two
		splitLines = append(
			splitLines,
			getSubstring(sourceLines[0], 0, width-*lastLineLength),
		)
		sourceLines[0] = getSubstring(sourceLines[0], width-*lastLineLength, len(sourceLines[0]))
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
		if i == 0 && (len(line)+*lastLineLength) < width {
			// Special case for the first line that has to takes into account currentLineLength
			splitLines[i] = line + strings.Repeat(" ", width-len(line))
		} else if i > 0 && len(line) < width {
			splitLines[i] = line + strings.Repeat(" ", width-len(line))
		}
	}

	for i, line := range splitLines {
		if i == 0 && len(*output) > 0 {
			(*output)[len(*output)-1] += stack.GetCurrent().Apply(line)
			*lastLineLength += len(line)
		} else if len(line) > 0 {
			// Then we decorate each line
			*output = append(*output, stack.GetCurrent().Apply(line))
			*lastLineLength = len(line)
		} else {
			*output = append(*output, "")
			*lastLineLength = 0
		}
	}
}
