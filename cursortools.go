package styledconsole

import (
	"fmt"
)

// hideCursor hides the cursor, can be reversed with showCursor()
func hideCursor() {
	fmt.Print("\033[?25l")
}

// showCursor restores the cursor after it was hidden
func showCursor() {
	fmt.Print("\033[?25h\033[?0c")
}

// clearWindowFromCursor clears all the output from the cursors' current position to the end of the screen.
func clearWindowFromCursor() {
	fmt.Print("\033[0J")
}
