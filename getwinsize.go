package styledconsole

import (
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

// GetWinsize return the size (width, height) of the current terminal window
func getWinsize() (int, int) {
	width, height, err := terminal.GetSize(int(os.Stdout.Fd()))

	if err != nil {
		return 120, 60
	}

	return width, height
}
