//go:build readline_notty
// +build readline_notty

package readline

import "errors"

func (rl *Instance) launchEditor(multiline []rune) ([]rune, error) {
	return rl.line.Runes(), errors.New("Not currently supported when compiled as `readline_notty`")
}
