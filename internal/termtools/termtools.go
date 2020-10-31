package termtools

import (
	"fmt"
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

type Position struct {
	X int
	Y int
}

// MoveUp moves the cursor a given number of lines up
func MoveCursorUp(lines int) {
	fmt.Printf("\033[%dA", lines)
}

// MoveDown moves the cursor a given number of lines down
func MoveCursorDown(lines int) {
	fmt.Printf("\033[%dB", lines)
}

// MoveRight moves the cursor a given number of lines to the right
func MoveCursorRight(columns int) {
	fmt.Printf("\033[%dC", columns)
}

// MoveLeft moves the cursor a given number of lines to the left
func MoveCursorLeft(columns int) {
	fmt.Printf("\033[%dD", columns)
}

// MoveToColumn moves the cursor a given column
func MoveCursorToColumn(column int) {
	fmt.Printf("\033[%dG", column)
}

// MoveToColumn moves the cursor a given row and column in the terminal
func MoveCursorToPosition(column int, row int) {
	fmt.Printf("\033[%d;%dH", row+1, column)
}

// SaveCursorPosition saves the position of the cursor, so it can be brought back with cursor.RestoreCursorPosition()
func SaveCursorPosition() {
	fmt.Print("\0337")
}

// RestoreCursorPosition brings back the cursor to the position it was saved with cursor.SaveCursorPosition()
func RestoreCursorPosition() {
	fmt.Print("\0338")
}

// HideCursor hides the cursor, can be reversed with cursor.ShowCursor()
func HideCursor() {
	fmt.Print("\033[?25l")
}

// ShowCursor restores the cursor after it was hidden
func ShowCursor() {
	fmt.Print("\033[?25h\033[?0c")
}

// ClearCursorLine clears all the output from the current line.
func ClearCursorLine() {
	fmt.Print("\033[2K")
}

// ClearCursorEndLine clears all the output from the current line after the current position.
func ClearCursorEndLine() {
	fmt.Print("\033[K")
}

// ClearWindowFromCursor clears all the output from the cursors' current position to the end of the screen.
func ClearWindowFromCursor() {
	fmt.Print("\033[0J")
}

// ClearWindow clears the entire screen.
func ClearWindow() {
	fmt.Print("\033[2J")
}

// GetCurrentPosition returns the current cursor position as a Position{X,Y} object
func GetCurrentCursorPosition() Position {
	if terminal.IsTerminal(int(os.Stdout.Fd())) {
		oldState, err := terminal.MakeRaw(int(os.Stdout.Fd()))
		if err != nil {
			return Position{X: 1, Y: 1}
		}

		// Make the terminal print the cursor position
		os.Stdout.Write([]byte("\033[6n"))
		var row, col int
		_, err = fmt.Scanf("\033[%d;%dR", &row, &col)
		if err != nil {
			return Position{X: 1, Y: 1}
		}

		err = terminal.Restore(int(os.Stdout.Fd()), oldState)
		if err != nil {
			return Position{X: 1, Y: 1}
		}

		return Position{X: col, Y: row}
	}

	return Position{X: 1, Y: 1}
}

// GetWinsize return the size (width, height) of the current terminal window
func GetWinsize() (int, int) {
	width, height, err := terminal.GetSize(int(os.Stdout.Fd()))

	if err != nil {
		return 120, 60
	}

	return width, height
}
