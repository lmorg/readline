package readline

func HkFnCursorMoveToStartOfLine(rl *Instance) {
	rl.viUndoSkipAppend = true
	if rl.line.RuneLen() == 0 {
		return
	}
	output := rl.clearHelpersStr()
	rl.line.SetCellPos(0)
	output += rl.echoStr()
	output += moveCursorForwardsStr(1)
	rl.print(output)
}

func HkFnCursorMoveToEndOfLine(rl *Instance) {
	rl.viUndoSkipAppend = true
	if rl.line.RuneLen() == 0 {
		return
	}
	output := rl.clearHelpersStr()
	rl.line.SetRunePos(rl.line.RuneLen())
	output += rl.echoStr()
	output += moveCursorForwardsStr(1)
	rl.print(output)
}

func HkFnClearAfterCursor(rl *Instance) {
	if rl.line.RuneLen() == 0 {
		return
	}
	output := rl.clearHelpersStr()
	rl.line.Set(rl, rl.line.Runes()[:rl.line.RunePos()])
	output += rl.echoStr()
	output += moveCursorForwardsStr(1)
	rl.print(output)
}

func HkFnClearLine(rl *Instance) {
	rl.clearPrompt()
	rl.resetHelpers()
}

func HkFnCursorJumpForwards(rl *Instance) {
	rl.viUndoSkipAppend = true
	output := rl.moveCursorByRuneAdjustStr(rl.viJumpE(tokeniseLine))
	rl.print(output)
}

func HkFnCursorJumpBackwards(rl *Instance) {
	rl.viUndoSkipAppend = true
	output := rl.moveCursorByRuneAdjustStr(rl.viJumpB(tokeniseLine))
	rl.print(output)
}

func _hkFnCancelActionModeAutoFind(rl *Instance) string {
	rl.viUndoSkipAppend = true
	output := rl.resetTabFindStr()
	output += rl.clearHelpersStr()
	rl.resetTabCompletion()
	output += rl.renderHelpersStr()
	return output
}
func _hkFnCancelActionModeTabFind(rl *Instance) string {
	rl.viUndoSkipAppend = true
	return rl.resetTabFindStr()
}
func _hkFnCancelActionModeViModeVimCommand(rl *Instance) string {
	rl.viUndoSkipAppend = true
	rl.vimCommandModeCancel()
	return rl.updateHelpersStr()
}
func _hkFnCancelActionModeTabCompletion(rl *Instance) string {
	rl.viUndoSkipAppend = true
	output := rl.clearHelpersStr()
	rl.resetTabCompletion()
	output += rl.renderHelpersStr()
	return output
}
func _hkFnCancelActionDefault(rl *Instance) string {
	rl.viUndoSkipAppend = true
	rl.modeViMode = vimKeys
	rl.viIteration = ""
	return rl.viHintMessageStr()
}

func HkFnRecallWord1(rl *Instance)    { hkFnRecallWord(rl, 1) }
func HkFnRecallWord2(rl *Instance)    { hkFnRecallWord(rl, 2) }
func HkFnRecallWord3(rl *Instance)    { hkFnRecallWord(rl, 3) }
func HkFnRecallWord4(rl *Instance)    { hkFnRecallWord(rl, 4) }
func HkFnRecallWord5(rl *Instance)    { hkFnRecallWord(rl, 5) }
func HkFnRecallWord6(rl *Instance)    { hkFnRecallWord(rl, 6) }
func HkFnRecallWord7(rl *Instance)    { hkFnRecallWord(rl, 7) }
func HkFnRecallWord8(rl *Instance)    { hkFnRecallWord(rl, 8) }
func HkFnRecallWord9(rl *Instance)    { hkFnRecallWord(rl, 9) }
func HkFnRecallWord10(rl *Instance)   { hkFnRecallWord(rl, 10) }
func HkFnRecallWord11(rl *Instance)   { hkFnRecallWord(rl, 11) }
func HkFnRecallWord12(rl *Instance)   { hkFnRecallWord(rl, 12) }
func HkFnRecallWordLast(rl *Instance) { hkFnRecallWord(rl, -1) }

func HkFnUndo(rl *Instance) {
	rl.viUndoSkipAppend = true
	if len(rl.viUndoHistory) == 0 {
		return
	}
	output := rl.undoLastStr()
	rl.viUndoSkipAppend = true
	rl.line.SetRunePos(rl.line.RuneLen())
	output += moveCursorForwardsStr(1)
	rl.print(output)
}

func HkFnClearScreen(rl *Instance) {
	if rl.isNoTty {
		return
	}

	rl.viUndoSkipAppend = true
	if rl.previewMode != previewModeClosed {
		HkFnModePreviewToggle(rl)
	}
	output := seqSetCursorPosTopLeft + seqClearScreen
	output += rl.echoStr()
	output += rl.renderHelpersStr()
	rl.print(output)
}

func HkFnModeFuzzyFind(rl *Instance) {
	if rl.isNoTty {
		return
	}

	rl.viUndoSkipAppend = true
	if !rl.modeTabCompletion {
		rl.modeAutoFind = true
		rl.getTabCompletion()
	}

	rl.modeTabFind = true
	rl.print(rl.updateTabFindStr([]rune{}))
}

func HkFnModeSearchHistory(rl *Instance) {
	if rl.isNoTty {
		return
	}

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
	rl.print(rl.updateTabFindStr([]rune{}))
}

func HkFnModeAutocomplete(rl *Instance) {
	if rl.isNoTty {
		return
	}

	rl.viUndoSkipAppend = true
	if rl.modeTabCompletion {
		rl.moveTabCompletionHighlight(1, 0)
	} else {
		rl.getTabCompletion()
	}

	if rl.previewMode == previewModeOpen || rl.previewRef == previewRefLine {
		rl.previewMode = previewModeAutocomplete
	}

	rl.print(rl.renderHelpersStr())
}

func HkFnCancelAction(rl *Instance) {
	switch {
	case rl.modeAutoFind && !rl.isNoTty:
		rl.print(_hkFnCancelActionModeAutoFind(rl))

	case rl.modeTabFind && !rl.isNoTty:
		rl.print(_hkFnCancelActionModeTabFind(rl))

	case rl.modeViMode == vimCommand:
		rl.print(_hkFnCancelActionModeViModeVimCommand(rl))

	case rl.modeTabCompletion && !rl.isNoTty:
		rl.print(_hkFnCancelActionModeTabCompletion(rl))

	default:
		rl.print(_hkFnCancelActionDefault(rl))
	}
}

func HkFnModePreviewToggle(rl *Instance) {
	if rl.isNoTty || rl.PreviewLine == nil {
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
	if rl.isNoTty {
		return
	}

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
		rl.print(rl.clearPreviewStr())

	case previewModeAutocomplete:
		rl.print(rl.clearPreviewStr())
		rl.resetHelpers()
	}

	output += rl.echoStr()
	output += rl.renderHelpersStr()
	rl.print(output)
}

func HkFnModePreviewLine(rl *Instance) {
	if rl.isNoTty || rl.PreviewLine == nil {
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
		rl.print(rl.renderHelpersStr())
	}
}
