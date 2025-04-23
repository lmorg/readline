//go:build readline_notty
// +build readline_notty

package readline

func HkFnClearScreen(rl *Instance)       {}
func HkFnModeFuzzyFind(rl *Instance)     {}
func HkFnModeSearchHistory(rl *Instance) {}
func HkFnModeAutocomplete(rl *Instance)  {}
func HkFnModePreviewToggle(rl *Instance) {}
func HkFnModePreviewLine(rl *Instance)   {}

func HkFnCancelAction(rl *Instance) {
	switch {
	case rl.modeViMode == vimCommand:
		print(_hkFnCancelActionModeViModeVimCommand(rl))

	default:
		print(_hkFnCancelActionDefault(rl))
	}
}
