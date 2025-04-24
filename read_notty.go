//go:build readline_notty
// +build readline_notty

package readline

import "errors"

func (rl *Instance) KeyPress(b []byte) {
	go func() {
		rl._noTtyKeyPress <- b
	}()
}

func read(p []byte) (int, error) {
	b, ok := <-rl._noTtyKeyPress

	if !ok {
		return 0, errors.New("channel closed")
	}

	copy(p, b)
	return len(b), nil
}
