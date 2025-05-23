package readline

import "os"

// GetTermWidth returns the width of Stdout or 80 if the width cannot be established
func GetTermWidth() (termWidth int) {
	var err error
	fd := int(os.Stdout.Fd())
	termWidth, _, err = GetSize(fd)
	if err != nil {
		termWidth = 80 // default to 80 with term width unknown as that is the de factor standard on older terms.
	}

	return
}

func (rl *Instance) termWidth() int {
	return rl._termWidth
}

func (rl *Instance) cacheTermWidth() {
	if rl.isNoTty {
		return
	}

	rl._termWidth = GetTermWidth()
}
