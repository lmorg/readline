package readline

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/mattn/go-runewidth"
)

func (rl *Instance) printf(format string, a ...interface{}) {
	s := fmt.Sprintf(format, a...)
	rl.print(s)
}

func (rl *Instance) print(s string) {
	if rl.isNoTty {
		return
	}

	_print(s)
}

func (rl *Instance) printErr(s string) {
	if rl.isNoTty {
		return
	}

	_printErr(s)
}

// rxAnsiEscape matches ANSI escape sequences that occupy zero visible
// columns on screen. Stripping these before measuring is required for
// any prompt or hint that uses richer output than plain SGR colors.
//
// Alternatives, in order:
//
//  1. OSC (Operating System Command): ESC ] ... (BEL | ESC \)
//     Hyperlinks (OSC 8), window titles (OSC 0/1/2), CWD reporting
//     (OSC 7), iTerm shell-integration marks (OSC 1337), notifications
//     (OSC 9), etc. Both string terminators (0x07 BEL and ESC \) are
//     accepted; the payload may not contain ESC or BEL.
//
//  2. CSI (Control Sequence Introducer): ESC [ params intermediates final
//     Covers SGR (final byte 'm'), cursor moves, clears, scroll regions
//     and friends. ECMA-48 conformant byte ranges: params 0x30-0x3F,
//     intermediates 0x20-0x2F, final 0x40-0x7E.
//
//  3. Single-byte Fe escapes: ESC X where X is 0x40-0x5F
//     Picks up DEC save/restore cursor (ESC 7 / ESC 8 are 0x37/0x38 so
//     handled separately below), index (ESC D), next line (ESC E),
//     reverse index (ESC M), single shifts, etc. Excludes the CSI ([)
//     and OSC (]) introducers, which are consumed by the earlier
//     alternatives via leftmost-first matching.
//
//  4. DEC private cursor save/restore: ESC 7, ESC 8
//     Numeric finals fall outside the Fe range above.
var rxAnsiEscape = regexp.MustCompile(
	`\x1b\][^\x07\x1b]*(?:\x07|\x1b\\)` +
		`|\x1b\[[\x30-\x3f]*[\x20-\x2f]*[\x40-\x7e]` +
		`|\x1b[\x40-\x5f]` +
		`|\x1b[78]`,
)

// strLen returns the number of terminal cells a string would occupy on a
// monospace display, after stripping zero-width ANSI escape sequences.
func strLen(s string) int {
	s = rxAnsiEscape.ReplaceAllString(s, "")
	return runewidth.StringWidth(s)
}

func (rl *Instance) echoStr() string {
	if len(rl.multiSplit) == 0 {
		rl.syntaxCompletion()
	}

	lineX, lineY := rl.lineWrapCellLen()
	posX, posY := rl.lineWrapCellPos()

	// reset cursor to start
	line := "\r"
	if posY > 0 {
		line += fmt.Sprintf(cursorUpf, posY)
	}

	// clear the line
	line += strings.Repeat("\x1b[2K\n", lineY+1) // clear line + move cursor down 1
	line += fmt.Sprintf(cursorUpf, lineY+1)
	//line += seqClearScreenBelow

	promptLen := rl.promptLen
	if promptLen < rl.termWidth() {
		line += rl.prompt
	} else {
		promptLen = 0
	}

	switch {
	case rl.PasswordMask != 0:
		line += strings.Repeat(string(rl.PasswordMask), rl.line.CellLen())

	case rl.line.CellLen()+promptLen > rl.termWidth():
		fallthrough

	case rl.SyntaxHighlighter == nil:
		line += strings.Join(lineWrap(rl, rl.termWidth()), "\r\n")

	default:
		syntax := rl.cacheSyntax.Get(rl.line.Runes())
		if len(syntax) == 0 {
			syntax = rl.SyntaxHighlighter(rl.line.Runes())

			if rl.DelayedSyntaxWorker == nil {
				rl.cacheSyntax.Append(rl.line.Runes(), syntax)
			}
		}
		line += syntax
	}

	y := lineY - posY
	if y > 0 {
		line += fmt.Sprintf(cursorUpf, y)
	}
	x := lineX - posX + 1
	if x > 0 {
		line += fmt.Sprintf(cursorBackf, x)
	}
	//print(line)
	return line
}

func lineWrap(rl *Instance, termWidth int) []string {
	var promptLen int
	if rl.promptLen < termWidth {
		promptLen = rl.promptLen
	}

	var (
		wrap       []string
		wrapRunes  [][]rune
		bufCellLen int
		length     = termWidth - promptLen
		line       = rl.line.Runes() //append(rl.line.Runes(), []rune{' ', ' '}...) // double space to work around wide characters
		lPos       int
	)

	wrapRunes = append(wrapRunes, []rune{})

	for r := range line {
		w := runewidth.RuneWidth(line[r])
		if bufCellLen+w > length {
			wrapRunes = append(wrapRunes, []rune(strings.Repeat(" ", promptLen)))
			lPos++
			bufCellLen = 0
		}
		bufCellLen += w
		wrapRunes[lPos] = append(wrapRunes[lPos], line[r])
	}

	wrap = make([]string, lPos+1)
	for i := range wrap {
		wrap[i] = string(wrapRunes[i])
	}

	return wrap
}

func (rl *Instance) lineWrapCellLen() (x, y int) {
	return LineWrappedCellPos(rl.promptLen, rl.line.Runes(), rl.termWidth())
}

func (rl *Instance) lineWrapCellPos() (x, y int) {
	return LineWrappedCellPos(rl.promptLen, rl.line.Runes()[:rl.line.RunePos()], rl.termWidth())
}

// LineWrappedCellPos is a unicode and wide character aware function for
// determining the x/y coordinates of a cell.
func LineWrappedCellPos(promptLen int, line []rune, termWidth int) (x, y int) {
	if promptLen >= termWidth {
		promptLen = 0
	}

	// avoid divide by zero error
	if termWidth-promptLen == 0 {
		return 0, 0
	}

	x = promptLen
	for i := range line {
		w := runewidth.RuneWidth(line[i])
		if x+w > termWidth {
			x = promptLen
			y++
		}
		x += w
	}

	return
}

func (rl *Instance) clearPrompt() {
	if rl.line.RuneLen() == 0 {
		return
	}

	output := rl.moveCursorToStartStr()

	if rl.termWidth() > rl.promptLen {
		output += strings.Repeat(" ", rl.termWidth()-rl.promptLen)
	}
	output += seqClearScreenBelow

	output += moveCursorBackwardsStr(rl.termWidth())
	output += rl.prompt

	rl.line.Set(rl, []rune{})
	rl.line.SetRunePos(0)

	rl.print(output)
}

func (rl *Instance) resetHelpers() {
	rl.modeAutoFind = false
	output := rl.clearPreviewStr()
	output += rl.clearHelpersStr()

	rl.resetHintText()
	rl.resetTabCompletion()

	rl.print(output)
}

func (rl *Instance) clearHelpersStr() string {
	posX, posY := rl.lineWrapCellPos()
	_, lineY := rl.lineWrapCellLen()

	y := lineY - posY

	output := moveCursorDownStr(y)
	output += "\r\n" + seqClearScreenBelow

	output += moveCursorUpStr(y + 1)
	output += moveCursorForwardsStr(posX)

	return output
}

func (rl *Instance) renderHelpersStr() string {
	output := rl.writeHintTextStr()
	output += rl.writeTabCompletionStr()
	output += rl.writePreviewStr()
	return output
}

func (rl *Instance) updateHelpersStr() string {
	rl.tcOffset = 0
	rl.getHintText()
	if rl.modeTabCompletion.Load() {
		rl.getTabCompletion()
	}
	output := rl.clearHelpersStr()
	output += rl.renderHelpersStr()

	return output
}
