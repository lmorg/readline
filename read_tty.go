//go:build !js && !readline_notty
// +build !js,!readline_notty

package readline

import (
	"os"
)

func (rl *Instance) KeyPress(b []byte) {
	panic("not supported without `readline_notty` build tag")
}

func read(b []byte) (int, error) {
	return os.Stdin.Read(b)
}
