//+build !windows

package terminal

import "golang.org/x/sys/unix"

const maxLineLength = 120

// GetWinsize return the size (width, height) of the current terminal window
func GetWinsize() (int, int) {
	ws, err := unix.IoctlGetWinsize(1, unix.TIOCGWINSZ) // file descriptor 1 is stdout

	if err != nil {
		return 60, 60
	}
	height := ws.Row
	width := ws.Col
	if width > maxLineLength {
		width = maxLineLength
	}

	return int(width), int(height)
}
