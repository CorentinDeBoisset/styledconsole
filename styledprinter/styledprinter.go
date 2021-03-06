// Package styledprinter helps to do some pretty-printing of text in terminal using ANSI escape sequences
package styledprinter

import (
	"fmt"
	"strings"
)

// WriteBlock prints a block of text using a string padding, with optionnal styles
func WriteBlock(message string, padding string, baseStyle string, newLine bool) {
	width, _ := getWinsize()

	widthWithoutPadding := width - len(padding)
	extractedBaseStyle := NewOutputStyle(baseStyle)

	// We ensure there is a last line at the end to have the background everywhere
	if message[len(message)-1:] != "\n" {
		message = message + "\n"
	}

	formattedLines := formatText(message, widthWithoutPadding, baseStyle)
	emptyLine := extractedBaseStyle.Apply(padding + strings.Repeat(" ", widthWithoutPadding))

	// However, we remove the last empty line, to blend with the block
	fmt.Printf("%s\n", emptyLine)
	for _, line := range formattedLines[0 : len(formattedLines)-1] {
		fmt.Printf("%s%s\n", extractedBaseStyle.Apply(padding), line)
	}
	fmt.Printf("%s\n", emptyLine)

	if newLine {
		fmt.Printf("\n")
	}
}

// Write prints a list of messages, one per line, with an optionnal end-of-line at the end
func Write(message string, newLine bool) {
	width, _ := getWinsize()

	formattedLines := formatText(message, width, "")
	fmt.Print(strings.Join(formattedLines, "\n"))
	if newLine {
		fmt.Print("\n")
	}
}
