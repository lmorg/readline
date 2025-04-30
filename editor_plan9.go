//go:build plan9
// +build plan9

package readline

import "errors"

func (rl *Instance) _launchEditor(multiline []rune) ([]rune, error) {
	return rl.line.Runes(), errors.New("Not currently supported on Plan 9")
}
