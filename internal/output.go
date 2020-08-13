package internal

import (
	"fmt"
	"strings"

	"github.com/corentindeboisset/styledconsole/internal/terminal"
)

// WriteBlock prints a block of text using a string padding, with optionnal styles
func WriteBlock() {
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
