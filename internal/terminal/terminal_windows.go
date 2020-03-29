package terminal

// GetWinsize return the size (width, height) of the current terminal window
func GetWinsize() (uint16, uint16) {
	// FIXME develop the compatibility for windows
	return 60, 60
}
