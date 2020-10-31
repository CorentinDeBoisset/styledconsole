// nolint:deadcode,unused
package main

import (
	"fmt"
	"strings"

	"github.com/corentindeboisset/styledconsole/internal"
	"github.com/corentindeboisset/styledconsole/internal/questionhlpr"
)

var progressStarted bool

// var currentProgress, maxProgress int

// Section function
func Section(title string) {
	titleLen := len(title)
	underline := strings.Repeat("=", titleLen)

	internal.Write(fmt.Sprintf("<fg=yellow;options=bold>%s\n%s\n</>", title, underline), true)
}

// Text function
func Text(content string) {
	internal.Write(content, true)
}

// Listing function
func Listing(items []string) {
	for _, item := range items {
		internal.Write(fmt.Sprintf(" <fg=yellow>*</> %s", item), true)
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

// ProgressStart function
func ProgressStart(totalSteps int) {
	if progressStarted {
		return
	}
	progressStarted = true
}

// ProgressAdvance function
func ProgressAdvance(stepCount int) {
	if !progressStarted {
		return
	}
}

// ProgressFinish function
func ProgressFinish() {
	if !progressStarted {
		return
	}
	progressStarted = false
}

// Ask function
func Ask(question string, validator func(string) bool) (string, error) {
	q := questionhlpr.Question{
		Label:         question,
		IsClosed:      false,
		IsHidden:      false,
		DefaultAnswer: "",
		Validator:     validator,
	}

	res, err := questionhlpr.AskQuestion(q)
	if err != nil {
		return "", err
	}
	return res, nil
}

// AskWithDefault function
func AskWithDefault(question string, defaultAnswer string, validator func(string) bool) (string, error) {
	q := questionhlpr.Question{
		Label:         question,
		IsClosed:      false,
		IsHidden:      false,
		DefaultAnswer: defaultAnswer,
		Validator:     validator,
	}

	res, err := questionhlpr.AskQuestion(q)
	if err != nil {
		return "", err
	}
	return res, nil
}

// AskHidden function
func AskHidden(question string, validator func(string) bool) (string, error) {
	q := questionhlpr.Question{
		Label:         question,
		IsClosed:      false,
		IsHidden:      true,
		DefaultAnswer: "",
		Validator:     validator,
	}

	res, err := questionhlpr.AskQuestion(q)
	if err != nil {
		return "", err
	}
	return res, nil
}

// Confirm function
func Confirm(question string) (bool, error) {
	return true, nil
}

// ConfirmWithDefault function
func ConfirmWithDefault(question string, defaultAnswer bool) (bool, error) {
	return true, nil
}

// Choice function
func Choice(question string, choices []string) (string, error) {
	q := questionhlpr.Question{
		Label:         question,
		IsClosed:      true,
		Choices:       choices,
		DefaultChoice: -1,
	}

	choice, err := questionhlpr.AskQuestion(q)
	if err != nil {
		fmt.Printf("error: %s", err)
		return "", err
	}

	return choice, nil
}

// ChoiceWithDefault function
func ChoiceWithDefault(question string, choices []string, defaultAnswer int) (string, error) {
	q := questionhlpr.Question{
		Label:         question,
		IsClosed:      true,
		Choices:       choices,
		DefaultChoice: defaultAnswer,
	}

	choice, err := questionhlpr.AskQuestion(q)
	if err != nil {
		fmt.Printf("error: %s", err)
		return "", err
	}

	return choice, nil
}

// Success function
func Success(content string) {
	internal.WriteBlock(fmt.Sprintf("Success:\n%s", content), "  ", "bg=green;fg=black", true)
}

// Warning function
func Warning(content string) {
	internal.WriteBlock(fmt.Sprintf("Warning:\n%s", content), "# ", "fg=yellow", true)
}

// Error function
func Error(content string) {
	internal.WriteBlock(fmt.Sprintf("Error:\n%s", content), "  ", "bg=red;fg=black", true)
}
