//go:build !readline_notty
// +build !readline_notty

package readline

func (rl *Instance) SetNoTtyCallback(callback chan *NoTtyCallbackT) {
	panic("readline needs to be compiled with `readline_notty` to support this")
}

func noTtyCallback(rl *Instance) {}
