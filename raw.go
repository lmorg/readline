package readline

// GetSize returns the dimensions of the given terminal.
func GetSize(fd int) (width, height int, err error) {
	return getSize(fd)
}
