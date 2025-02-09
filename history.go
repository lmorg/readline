package readline

import (
	"errors"
	"strconv"
	"strings"
)

// History is an interface to allow you to write your own history logging
// tools. eg sqlite backend instead of a file system.
// By default readline will just use the dummyLineHistory interface which only
// logs the history to memory ([]string to be precise).
type History interface {
	// Append takes the line and returns an updated number of lines or an error
	Write(string) (int, error)

	// GetLine takes the historic line number and returns the line or an error
	GetLine(int) (string, error)

	// Len returns the number of history lines
	Len() int

	// Dump returns everything in readline. The return is an interface{} because
	// not all LineHistory implementations will want to structure the history in
	// the same way. And since Dump() is not actually used by the readline API
	// internally, this methods return can be structured in whichever way is most
	// convenient for your own applications (or even just create an empty
	//function which returns `nil` if you don't require Dump() either)
	Dump() interface{}
}

// ExampleHistory is an example of a LineHistory interface:
type ExampleHistory struct {
	items []string
}

// Write to history
func (h *ExampleHistory) Write(s string) (int, error) {
	h.items = append(h.items, s)
	return len(h.items), nil
}

// GetLine returns a line from history
func (h *ExampleHistory) GetLine(i int) (string, error) {
	switch {
	case i < 0:
		return "", errors.New("requested history item out of bounds: < 0")
	case i > h.Len()-1:
		return "", errors.New("requested history item out of bounds: > Len()")
	default:
		return h.items[i], nil
	}
}

// Len returns the number of lines in history
func (h *ExampleHistory) Len() int {
	return len(h.items)
}

// Dump returns the entire history
func (h *ExampleHistory) Dump() interface{} {
	return h.items
}

// NullHistory is a null History interface for when you don't want line
// entries remembered eg password input.
type NullHistory struct{}

// Write to history
func (h *NullHistory) Write(s string) (int, error) {
	return 0, nil
}

// GetLine returns a line from history
func (h *NullHistory) GetLine(i int) (string, error) {
	return "", nil
}

// Len returns the number of lines in history
func (h *NullHistory) Len() int {
	return 0
}

// Dump returns the entire history
func (h *NullHistory) Dump() interface{} {
	return []string{}
}

// Browse historic lines
func (rl *Instance) walkHistory(i int) {
	if rl.previewRef == previewRefLine {
		return // don't walk history if [f9] preview open
	}

	line := rl.line.String()
	line = strings.TrimSpace(line)
	rl._walkHistory(i, line)
}

func (rl *Instance) _walkHistory(i int, oldLine string) {
	var (
		newLine string
		err     error
	)

	switch rl.histPos + i {
	case -1, rl.History.Len() + 1:
		return

	case rl.History.Len():
		rl.clearPrompt()
		rl.histPos += i
		if len(rl.viUndoHistory) > 0 && rl.viUndoHistory[len(rl.viUndoHistory)-1].String() != "" {
			rl.line = rl.lineBuf.Duplicate()
		}

	default:
		newLine, err = rl.History.GetLine(rl.histPos + i)
		if err != nil {
			rl.resetHelpers()
			print("\r\n" + err.Error() + "\r\n")
			print(rl.prompt)
			return
		}

		if rl.histPos-i == rl.History.Len() {
			rl.lineBuf = rl.line.Duplicate()
		}

		rl.histPos += i
		if oldLine == newLine {
			rl._walkHistory(i, newLine)
			return
		}

		if len(rl.viUndoHistory) > 0 {
			last := rl.viUndoHistory[len(rl.viUndoHistory)-1]
			if !strings.HasPrefix(newLine, strings.TrimSpace(last.String())) {
				rl._walkHistory(i, oldLine)
				return
			}
		}

		rl.clearPrompt()

		rl.line = new(UnicodeT)
		rl.line.Set(rl, []rune(newLine))
	}

	if i > 0 {
		_, y := rl.lineWrapCellLen()
		print(strings.Repeat("\r\n", y))
		rl.line.SetRunePos(rl.line.RuneLen())
	} else {
		rl.line.SetCellPos(rl.termWidth - rl.promptLen - 1)
	}
	print(rl.echoStr())
	print(rl.updateHelpersStr())
}

func (rl *Instance) autocompleteHistory() ([]string, map[string]string) {
	if rl.AutocompleteHistory != nil {
		rl.tcPrefix = rl.line.String()
		return rl.AutocompleteHistory(rl.tcPrefix)
	}

	var (
		items []string
		descs = make(map[string]string)

		line string
		num  string
		err  error
	)

	rl.tcPrefix = rl.line.String()
	for i := rl.History.Len() - 1; i >= 0; i-- {
		line, err = rl.History.GetLine(i)
		if err != nil {
			continue
		}

		if !strings.HasPrefix(line, rl.tcPrefix) {
			continue
		}

		line = strings.Replace(line, "\n", ` `, -1)[rl.line.RuneLen():]

		if descs[line] != "" {
			continue
		}

		items = append(items, line)
		num = strconv.Itoa(i)

		descs[line] = num
	}

	return items, descs
}
