// Package styledconsole helps to make your GUI tools user-friendly with methods to pretty-print text, lists, tables, user prompts, progress bars...
package styledconsole

import (
	"fmt"
	"strings"

	"github.com/corentindeboisset/styledconsole/styledprinter"
)

// Section displays the given string as the title of some command section.
func Section(title string) {
	titleLen := len(title)
	underline := strings.Repeat("=", titleLen)

	styledprinter.Write(fmt.Sprintf("<fg=yellow;options=bold>%s\n%s\n</>", title, underline), true)
}

// Text displays the given string as regular text. This is useful to render help messages and instructions for the user running the command.
// This methods support style tags such as "<fg=blue>blue text</>".
func Text(content string) {
	styledprinter.Write(content, true)
}

// Listing displays an list of elements
func Listing(items []string) {
	for _, item := range items {
		styledprinter.Write(fmt.Sprintf(" <fg=yellow>*</> %s", item), true)
	}
}

// Table pretty-prints a table with headers. It does not support multiline cells or sytling.
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

// NewLine prints a line break.
func NewLine() {
	styledprinter.Write("", true)
}

// NewLines print the given amount of new breaks.
func NewLines(newLineCount int) {
	if newLineCount > 0 {
		styledprinter.Write(strings.Repeat("\n", newLineCount-1), true)
	}
}

// Ask prompts a question with the given label.
// A function can be given to ensure the validity of the response. To allow any response (even empty), put nil as validator.
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

// Same as Ask() but if the user's answer is empty, the given defaultAnswer is chosen instead.
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

// Same as Ask() but the characters typed by the user are not printed in the output, in a linux-style password prompt.
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

// Confirm prompts a yes/no question.
func Confirm(label string) (bool, error) {
	return askConfirm(label, nil)
}

// ConfirmWithDefault prompts a yes/no question, with a given answer by default if the user's answer is empty.
func ConfirmWithDefault(label string, defaultAnswer bool) (bool, error) {
	return askConfirm(label, &defaultAnswer)
}

// Choice prints a list of choices the user can choose between.
// The prompts adapts itself to the size of the terminal.
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

// ChoiceWithDefault is the same as Choice() but a specific answer index should be given to highlight by default.
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

// Success displays the given string highlighted as a successful message (with a green background and an [OK] label).
func Success(content string) {
	styledprinter.WriteBlock(fmt.Sprintf("Success:\n%s", content), "  ", "bg=green;fg=black", true)
}

// Warning displays the given string highlighted as a warning message (with yellow text and a [Warning] label).
func Warning(content string) {
	styledprinter.WriteBlock(fmt.Sprintf("Warning:\n%s", content), "# ", "fg=yellow", true)
}

// Error displays the given string highlighted as an error message (with a red background and the [Error] label).
func Error(content string) {
	styledprinter.WriteBlock(fmt.Sprintf("Error:\n%s", content), "  ", "bg=red;fg=black", true)
}
