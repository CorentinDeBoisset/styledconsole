// Helper functions to pretty-print messages in a terminal.
package styledconsole

import (
	"fmt"
	"strings"

	"github.com/corentindeboisset/styledconsole/styledprinter"
)

// Section prints a section title
func Section(title string) {
	titleLen := len(title)
	underline := strings.Repeat("=", titleLen)

	styledprinter.Write(fmt.Sprintf("<fg=yellow;options=bold>%s\n%s\n</>", title, underline), true)
}

// Text function
func Text(content string) {
	styledprinter.Write(content, true)
}

// Listing function
func Listing(items []string) {
	for _, item := range items {
		styledprinter.Write(fmt.Sprintf(" <fg=yellow>*</> %s", item), true)
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
	styledprinter.Write("", true)
}

// NewLines function
func NewLines(newLineCount int) {
	if newLineCount > 0 {
		styledprinter.Write(strings.Repeat("\n", newLineCount-1), true)
	}
}

// Prompts a question with the given label.
// A function can be given to ensure the validity of the response. To allow any response (even empty), put nil as validator
func Ask(label string, validator func(string) bool) (string, error) {
	q := question{
		Label:         label,
		IsClosed:      false,
		IsHidden:      false,
		DefaultAnswer: "",
		Validator:     validator,
	}

	res, err := askQuestion(q)
	if err != nil {
		return "", err
	}
	return res, nil
}

// Same as Ask() but the characters typed by the user are not printed in the stdout, in a linux-style password prompt
func AskWithDefault(label string, defaultAnswer string, validator func(string) bool) (string, error) {
	q := question{
		Label:         label,
		IsClosed:      false,
		IsHidden:      false,
		DefaultAnswer: defaultAnswer,
		Validator:     validator,
	}

	res, err := askQuestion(q)
	if err != nil {
		return "", err
	}
	return res, nil
}

// AskHidden function
func AskHidden(label string, validator func(string) bool) (string, error) {
	q := question{
		Label:         label,
		IsClosed:      false,
		IsHidden:      true,
		DefaultAnswer: "",
		Validator:     validator,
	}

	res, err := askQuestion(q)
	if err != nil {
		return "", err
	}
	return res, nil
}

// Confirm function
func Confirm(label string) (bool, error) {
	return askConfirm(label, nil)
}

// ConfirmWithDefault function
func ConfirmWithDefault(label string, defaultAnswer bool) (bool, error) {
	return askConfirm(label, &defaultAnswer)
}

// Choice function
func Choice(label string, choices []string) (string, error) {
	q := question{
		Label:         label,
		IsClosed:      true,
		Choices:       choices,
		DefaultChoice: -1,
	}

	choice, err := askQuestion(q)
	if err != nil {
		fmt.Printf("error: %s", err)
		return "", err
	}

	return choice, nil
}

// ChoiceWithDefault function
func ChoiceWithDefault(label string, choices []string, defaultAnswer int) (string, error) {
	q := question{
		Label:         label,
		IsClosed:      true,
		Choices:       choices,
		DefaultChoice: defaultAnswer,
	}

	choice, err := askQuestion(q)
	if err != nil {
		fmt.Printf("error: %s", err)
		return "", err
	}

	return choice, nil
}

// Success function
func Success(content string) {
	styledprinter.WriteBlock(fmt.Sprintf("Success:\n%s", content), "  ", "bg=green;fg=black", true)
}

// Warning function
func Warning(content string) {
	styledprinter.WriteBlock(fmt.Sprintf("Warning:\n%s", content), "# ", "fg=yellow", true)
}

// Error function
func Error(content string) {
	styledprinter.WriteBlock(fmt.Sprintf("Error:\n%s", content), "  ", "bg=red;fg=black", true)
}
