//go:build windows || js || plan9 || readline_notty
// +build windows js plan9 readline_notty

package readline

func (rl *Instance) sigwinch() {
	rl.closeSigwinch = func() {
		// empty function because SIGWINCH isn't supported on these platforms
	}
}
