// nolint:deadcode,unused
package main

import (
	"fmt"
	"strings"

	"github.com/corentindeboisset/styledconsole/internal"
)

// Section function
func Section(title string) {
	titleLen := len(title)
	underline := strings.Repeat("=", titleLen)

	internal.Write(fmt.Sprintf("<fg=yellow>%s\n%s</>", title, underline), true)
}

// Text function
func Text(content string) {
	internal.Write(content, true)
}

// Listing function
func Listing(items []string) {
	for _, item := range items {
		internal.Write(fmt.Sprintf(" * %s", item), true)
	}
}

// Table pretty-prints a table with headers. It does not support multiline cells or sytling
func Table(headers []string, rows [][]string) {
	// First we have to determinate the width of every column
	var columnWidths []int
	for _, headerItem := range headers {
		columnWidths = append(columnWidths, len(headerItem))
	}
	for _, row := range rows {
		for columnNb, rowItem := range row {
			if columnNb < len(columnWidths) {
				if columnWidths[columnNb] < len(rowItem) {
					columnWidths[columnNb] = len(rowItem)
				}
			} else {
				columnWidths = append(columnWidths, len(rowItem))
			}
		}
	}

	// Prepare the row spearator
	sectionSeparator := "+"
	for _, width := range columnWidths {
		sectionSeparator += fmt.Sprintf("%s+", strings.Repeat("-", width+2))
	}

	var lines []string
	lines = append(lines, sectionSeparator)
	rowToPrint := "|"
	for columnNb, headerItem := range headers {
		rowToPrint += fmt.Sprintf(" %s%s |", headerItem, strings.Repeat(" ", columnWidths[columnNb]-len(headerItem)))
	}
	lines = append(lines, rowToPrint)
	lines = append(lines, sectionSeparator)

	for _, row := range rows {
		rowToPrint = "|"
		for columnNb, rowItem := range row {
			rowToPrint += fmt.Sprintf(" %s%s |", rowItem, strings.Repeat(" ", columnWidths[columnNb]-len(rowItem)))
		}
		lines = append(lines, rowToPrint)
	}
	lines = append(lines, sectionSeparator)

	fmt.Printf("%s\n", strings.Join(lines, "\n"))
}

// NewLine function
func NewLine() {
	internal.Write("", true)
}

// NewLines function
func NewLines(newLineCount int) {
	if newLineCount > 0 {
		internal.Write(strings.Repeat("\n", newLineCount-1), true)
	}
}

// Note function
func Note(content string) {
}

// Caution function
func Caution(content string) {
}

// ProgressStart function
func ProgressStart(totalSteps int) {
}

// ProgressAdvance function
func ProgressAdvance(stepCount int) {
}

// ProgressFinish function
func ProgressFinish() {
}

// Ask function
func Ask(question string, validator func(string) bool) string {
	return ""
}

// AskWithDefault function
func AskWithDefault(question string, defaultAnswer string, validator func(string) bool) string {
	return ""
}

// AskHidden function
func AskHidden(question string, validator func(string) bool) string {
	return ""
}

// Confirm function
func Confirm(question string) bool {
	return true
}

// ConfirmWithDefault function
func ConfirmWithDefault(question string, defaultAnswer bool) bool {
	return true
}

// Choice function
func Choice(question string, choices []string) int {
	return 0
}

// ChoiceWithDefault function
func ChoiceWithDefault(question string, choices []string, defaultAnswer int) int {
	return 0
}

// Success function
func Success(content string) {
	internal.WriteBlock(fmt.Sprintf("Error:\n%s", content), "  ", "bg=green;fg=white", true)
}

// Warning function
func Warning(content string) {
	internal.WriteBlock(fmt.Sprintf("Error:\n%s", content), "  ", "bg=red;fg=white", true)
}

// Error function
func Error(content string) {
	internal.WriteBlock(fmt.Sprintf("Error:\n%s", content), "  ", "bg=red;fg=white", true)
}
