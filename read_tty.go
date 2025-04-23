//go:build !js && !readline_notty
// +build !js,!readline_notty

package readline

import (
	"os"
)

func read(b []byte) (int, error) {
	return os.Stdin.Read(b)
}
