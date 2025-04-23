//go:build !readline_notty
// +build !readline_notty

package readline

func HkFnClearScreen(rl *Instance) {
	rl.viUndoSkipAppend = true
	if rl.previewMode != previewModeClosed {
		HkFnModePreviewToggle(rl)
	}
	output := seqSetCursorPosTopLeft + seqClearScreen
	output += rl.echoStr()
	output += rl.renderHelpersStr()
	print(output)
}

func HkFnModeFuzzyFind(rl *Instance) {
	rl.viUndoSkipAppend = true
	if !rl.modeTabCompletion {
		rl.modeAutoFind = true
		rl.getTabCompletion()
	}

	rl.modeTabFind = true
	print(rl.updateTabFindStr([]rune{}))
}

func HkFnModeSearchHistory(rl *Instance) {
	rl.viUndoSkipAppend = true
	rl.modeAutoFind = true
	rl.tcOffset = 0
	rl.modeTabCompletion = true
	rl.tcDisplayType = TabDisplayMap
	rl.tabMutex.Lock()
	rl.tcSuggestions, rl.tcDescriptions = rl.autocompleteHistory()
	rl.tabMutex.Unlock()
	rl.initTabCompletion()

	rl.modeTabFind = true
	print(rl.updateTabFindStr([]rune{}))
}

func HkFnModeAutocomplete(rl *Instance) {
	rl.viUndoSkipAppend = true
	if rl.modeTabCompletion {
		rl.moveTabCompletionHighlight(1, 0)
	} else {
		rl.getTabCompletion()
	}

	if rl.previewMode == previewModeOpen || rl.previewRef == previewRefLine {
		rl.previewMode = previewModeAutocomplete
	}

	print(rl.renderHelpersStr())
}

func HkFnCancelAction(rl *Instance) {
	switch {
	case rl.modeAutoFind:
		print(_hkFnCancelActionModeAutoFind(rl))

	case rl.modeTabFind:
		print(_hkFnCancelActionModeTabFind(rl))

	case rl.modeViMode == vimCommand:
		print(_hkFnCancelActionModeViModeVimCommand(rl))

	case rl.modeTabCompletion:
		print(_hkFnCancelActionModeTabCompletion(rl))

	default:
		print(_hkFnCancelActionDefault(rl))
	}
}

func HkFnModePreviewToggle(rl *Instance) {
	if rl.PreviewLine == nil {
		return
	}
	if !rl.modeAutoFind && !rl.modeTabCompletion && !rl.modeTabFind &&
		rl.previewMode == previewModeClosed {

		if rl.modeTabCompletion {
			rl.moveTabCompletionHighlight(1, 0)
		} else {
			rl.getTabCompletion()
		}
		defer func() { rl.previewMode++ }()
	}

	_fnPreviewToggle(rl)
}

func _fnPreviewToggle(rl *Instance) {
	rl.viUndoSkipAppend = true
	var output string

	switch rl.previewMode {
	case previewModeClosed:
		output = curPosSave + seqSaveBuffer + seqClearScreen
		rl.previewMode++
		size, _ := rl.getPreviewXY()
		if size != nil {
			output += rl.previewMoveToPromptStr(size)
		}

	case previewModeOpen:
		print(rl.clearPreviewStr())

	case previewModeAutocomplete:
		print(rl.clearPreviewStr())
		rl.resetHelpers()
	}

	output += rl.echoStr()
	output += rl.renderHelpersStr()
	print(output)
}

func HkFnModePreviewLine(rl *Instance) {
	if rl.PreviewLine == nil {
		return
	}
	if rl.PreviewInit != nil {
		// forced rerun of command line preview
		rl.PreviewInit()
		rl.previewCache = nil
	}

	if !rl.modeAutoFind && !rl.modeTabCompletion && !rl.modeTabFind &&
		rl.previewMode == previewModeClosed {
		defer func() { rl.previewMode++ }()
	}

	rl.previewRef = previewRefLine

	if rl.previewMode == previewModeClosed {
		_fnPreviewToggle(rl)
	} else {
		print(rl.renderHelpersStr())
	}
}
