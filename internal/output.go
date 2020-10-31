package internal

import (
	"fmt"
	"strings"

	"github.com/corentindeboisset/styledconsole/internal/style"
	"github.com/corentindeboisset/styledconsole/internal/termtools"
)

// WriteBlock prints a block of text using a string padding, with optionnal styles
func WriteBlock(message string, padding string, baseStyle string, newLine bool) {
	width, _ := termtools.GetWinsize()

	widthWithoutPadding := width - len(padding)
	extractedBaseStyle := style.NewOutputStyle(baseStyle)

	// We ensure there is a last line at the end to have the background everywhere
	if message[len(message)-1:] != "\n" {
		message = message + "\n"
	}

	formattedLines := FormatText(message, widthWithoutPadding, extractedBaseStyle)
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
	width, _ := termtools.GetWinsize()

	formattedLines := FormatText(message, width, nil)
	fmt.Print(strings.Join(formattedLines, "\n"))
	if newLine {
		fmt.Print("\n")
	}
}
