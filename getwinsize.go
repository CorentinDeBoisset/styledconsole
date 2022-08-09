package styledconsole

import (
	"os"

	"golang.org/x/term"
)

// GetWinsize return the size (width, height) of the current terminal window
func getWinsize() (int, int) {
	width, height, err := term.GetSize(int(os.Stdout.Fd()))

	if err != nil {
		return 120, 60
	}

	return width, height
}
