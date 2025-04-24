//go:build readline_notty
// +build readline_notty

package readline

func (rl *Instance) SetNoTtyCallback(callback chan *NoTtyCallbackT) {
	rl._noTtyKeyPress = make(chan []byte)
	rl._noTtyCallback = callback
}

func noTtyCallback(rl *Instance) {
	rl._noTtyCallback <- &NoTtyCallbackT{
		Line: rl.line.Duplicate(),
		Hint: string(rl.hintText),
	}
}

func (rl *Instance) Close() {
	close(rl._noTtyKeyPress)
	close(rl._noTtyCallback)
}
