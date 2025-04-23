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
	print(output)
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
	print(output)
}

func HkFnClearAfterCursor(rl *Instance) {
	if rl.line.RuneLen() == 0 {
		return
	}
	output := rl.clearHelpersStr()
	rl.line.Set(rl, rl.line.Runes()[:rl.line.RunePos()])
	output += rl.echoStr()
	output += moveCursorForwardsStr(1)
	print(output)
}

func HkFnClearLine(rl *Instance) {
	rl.clearPrompt()
	rl.resetHelpers()
}

func HkFnCursorJumpForwards(rl *Instance) {
	rl.viUndoSkipAppend = true
	output := rl.moveCursorByRuneAdjustStr(rl.viJumpE(tokeniseLine))
	print(output)
}

func HkFnCursorJumpBackwards(rl *Instance) {
	rl.viUndoSkipAppend = true
	output := rl.moveCursorByRuneAdjustStr(rl.viJumpB(tokeniseLine))
	print(output)
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
	print(output)
}
