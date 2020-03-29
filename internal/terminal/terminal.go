//+build !windows

package terminal

import "golang.org/x/sys/unix"

const defaultWidth = 60
const defaultHeight = 60
const maxLineLength = 120

// GetWinsize return the size (width, height) of the current terminal window
func GetWinsize() (uint16, uint16) {
	ws, err := unix.IoctlGetWinsize(1, unix.TIOCGWINSZ) // file descriptor 1 is stdout

	if err != nil {
		return 60, 60
	}
	height := ws.Row
	width := ws.Col
	if width > maxLineLength {
		width = maxLineLength
	}

	return width, height
}
