//go:build !windows && !js && !plan9
// +build !windows,!js,!plan9

package readline

import (
	"os"
	"os/signal"
	"syscall"
)

func (rl *Instance) sigwinch() {
	if rl.isNoTty {
		rl.closeSigwinch = func() {}
		return
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGWINCH)
	go func() {
		for range ch {

			width := GetTermWidth()

			switch {
			case !rl.modeTabCompletion || width == rl.termWidth():
				// no nothing

			case width < rl.termWidth():
				rl._termWidth = width
				HkFnClearScreen(rl)

			default:
				rl._termWidth = width
			}

		}
	}()

	rl.closeSigwinch = func() {
		signal.Stop(ch)
		close(ch)
	}
}
