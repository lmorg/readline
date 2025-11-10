package readline

import (
	"fmt"
	"strings"

	"github.com/lmorg/readline/v4/find"
)

func (rl *Instance) backspaceTabFindStr() string {
	if len(rl.tfLine) > 0 {
		rl.tfLine = rl.tfLine[:len(rl.tfLine)-1]
	}
	return rl.updateTabFindStr([]rune{})
}

func _updateTabFindHelpersStr(rl *Instance) (output string) {
	rl.tabMutex.Unlock()
	output = rl.clearHelpersStr()
	rl.initTabCompletion()
	output += rl.renderHelpersStr()
	return
}

func hintTextFindSearchStr(s string) []rune { return []rune(fmt.Sprintf("%s match: ", s)) }
func hintTextFindCancelStr(s string) []rune { return []rune(fmt.Sprintf("Cancelled %s match", s)) }

func (rl *Instance) updateTabFindStr(r []rune) string {
	rl.tfLine = append(rl.tfLine, r...)

	rl.tabMutex.Lock()

	if len(rl.tfLine) == 0 {
		rl.hintText = hintTextFindSearchStr("partial word")
		rl.tfSuggestions = append(rl.tcSuggestions, []string{}...)
		return _updateTabFindHelpersStr(rl)
	}

	find, err := find.New(string(rl.tfLine))
	rl.rFindSearch = hintTextFindSearchStr(find.Description())
	rl.rFindCancel = hintTextFindCancelStr(find.Description())
	if err != nil {
		rl.tfSuggestions = []string{err.Error()}
		return _updateTabFindHelpersStr(rl)
	}

	rl.hintText = append(rl.rFindSearch, rl.tfLine...)
	rl.hintText = append(rl.hintText, []rune(seqReset+seqBlink+"_"+seqReset)...)

	rl.tfSuggestions = make([]string, 0)
	for i := range rl.tcSuggestions {
		if find.MatchString(strings.TrimSpace(rl.tcSuggestions[i])) {
			rl.tfSuggestions = append(rl.tfSuggestions, rl.tcSuggestions[i])

		} else if rl.tcDisplayType == TabDisplayList && find.MatchString(rl.tcDescriptions[rl.tcSuggestions[i]]) {
			// this is a list so lets also check the descriptions
			rl.tfSuggestions = append(rl.tfSuggestions, rl.tcSuggestions[i])
		}
	}

	return _updateTabFindHelpersStr(rl)
}

func (rl *Instance) resetTabFindStr() string {
	rl.modeTabFind = false
	rl.tfLine = []rune{}
	if rl.modeAutoFind {
		rl.hintText = []rune{}
	} else {
		rl.hintText = rl.rFindCancel
	}
	rl.modeAutoFind = false

	output := rl.clearHelpersStr()
	rl.initTabCompletion()
	output += rl.renderHelpersStr()
	return output
}
