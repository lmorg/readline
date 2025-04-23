//go:build js && !readline_notty
// +build js,!readline_notty

package readline

import "errors"

func (rl *Instance) launchEditor(multiline []rune) ([]rune, error) {
	return rl.line.Runes(), errors.New("Not currently supported in WebAssembly")
}
