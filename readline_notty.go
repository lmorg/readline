//go:build readline_notty
// +build readline_notty

package readline

func (rl *Instance) SetNoTtyCallback(callback chan *NoTtyCallbackT) {
	rl.noTtyCallback = callback
}

func noTtyCallback(rl *Instance) {
	rl.noTtyCallback <- &NoTtyCallbackT{
		Line: rl.line.Duplicate(),
		Hint: string(rl.hintText),
	}
}
