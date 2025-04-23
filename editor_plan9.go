//go:build plan9 && !readline_notty
// +build plan9,!readline_notty

package readline

import "errors"

func (rl *Instance) launchEditor(multiline []rune) ([]rune, error) {
	return rl.line.Runes(), errors.New("Not currently supported on Plan 9")
}
