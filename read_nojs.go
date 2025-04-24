//go:build !js
// +build !js

package readline

import (
	"errors"
	"os"
)

func (rl *Instance) read(p []byte) (int, error) {
	if rl.isNoTty {
		return rl._readNoTty(p)
	}

	return os.Stdin.Read(p)
}

func (rl *Instance) _readNoTty(p []byte) (int, error) {
	b, ok := <-rl._noTtyKeyPress

	if !ok {
		return 0, errors.New("channel closed")
	}

	copy(p, b)
	return len(b), nil
}
