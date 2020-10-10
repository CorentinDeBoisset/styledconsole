package cursor

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
func MoveUp(lines int) {
	fmt.Printf("\033[%dA", lines)
}

// MoveDown moves the cursor a given number of lines down
func MoveDown(lines int) {
	fmt.Printf("\033[%dB", lines)
}

// MoveRight moves the cursor a given number of lines to the right
func MoveRight(columns int) {
	fmt.Printf("\033[%dC", columns)
}

// MoveLeft moves the cursor a given number of lines to the left
func MoveLeft(columns int) {
	fmt.Printf("\033[%dD", columns)
}

// MoveToColumn moves the cursor a given column
func MoveToColumn(column int) {
	fmt.Printf("\033[%dG", column)
}

// MoveToColumn moves the cursor a given row and column in the terminal
func MoveToPosition(column int, row int) {
	fmt.Printf("\033[%d;%dH", row+1, column)
}

// SavePosition saves the position of the cursor, so it can be brought back with cursor.RestorePosition()
func SavePosition() {
	fmt.Print("\0337")
}

// RestorePosition brings back the cursor to the position it was saved with cursor.SavePosition()
func RestorePosition() {
	fmt.Print("\0338")
}

// Hide hides the cursor, can be reversed with cursor.Show()
func Hide() {
	fmt.Print("\033[?25l")
}

// Show restores the cursor after it was hidden
func Show() {
	fmt.Print("\033[?25h\033[?0c")
}

// ClearLine clears all the output from the current line.
func ClearLine() {
	fmt.Print("\033[2K")
}

// ClearLineAfter clears all the output from the current line after the current position.
func ClearLineAfter() {
	fmt.Print("\033[K")
}

// ClearOutput clears all the output from the cursors' current position to the end of the screen.
func ClearOutput() {
	fmt.Print("\033[0J")
}

// ClearScreen clears the entire screen.
func ClearScreen() {
	fmt.Print("\033[2J")
}

// GetCurrentPosition returns the current cursor position as a Position{X,Y} object
func GetCurrentPosition() Position {
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
