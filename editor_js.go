//go:build js
// +build js

package readline

import "errors"

func (rl *Instance) _launchEditor(multiline []rune) ([]rune, error) {
	return rl.line.Runes(), errors.New("Not currently supported in WebAssembly")
}
