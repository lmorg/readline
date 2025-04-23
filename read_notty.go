//go:build readline_notty
// +build readline_notty

package readline

import "errors"

var _noTtyKeyPress = make(chan []byte)

func KeyPress(b []byte) {
	go func() {
		_noTtyKeyPress <- b
	}()
}

func read(p []byte) (int, error) {
	b, ok := <-_noTtyKeyPress

	if !ok {
		return 0, errors.New("channel closed")
	}

	copy(p, b)
	return len(b), nil
}
