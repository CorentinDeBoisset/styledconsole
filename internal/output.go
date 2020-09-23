package internal

import (
	"fmt"
	"strings"

	"github.com/corentindeboisset/styledconsole/internal/terminal"
)

// WriteBlock prints a block of text using a string padding, with optionnal styles
func WriteBlock(message string, padding string, baseStyle string, newLine bool) {
	width, _ := terminal.GetWinsize()

	widthWithoutPadding := width - len(padding)

	formattedLines := FormatText(message, widthWithoutPadding, baseStyle)
	for _, line := range formattedLines {
		fmt.Printf("%s%s\n", padding, line)
	}
	if newLine {
		fmt.Print("\n")
	}
}

// Write prints a list of messages, one per line, with an optionnal end-of-line at the end
func Write(message string, newLine bool) {
	width, _ := terminal.GetWinsize()

	currentLines := FormatText(message, width, "")
	fmt.Print(strings.Join(currentLines, "\n"))
	if newLine {
		fmt.Print("\n")
	}
}
