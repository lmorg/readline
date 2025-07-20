package readline

import (
	"fmt"
	"strconv"

	"github.com/mattn/go-runewidth"
)

func (rl *Instance) initTabGrid() {
	rl.tabMutex.Lock()
	defer rl.tabMutex.Unlock()

	var suggestions *suggestionsT
	if rl.modeTabFind {
		suggestions = newSuggestionsT(rl, rl.tfSuggestions)
	} else {
		suggestions = newSuggestionsT(rl, rl.tcSuggestions)
	}

	rl.tcMaxLength = rl.MinTabItemLength

	for i := 0; i < suggestions.Len(); i++ {
		l := suggestions.ItemLen(i)
		if l > rl.tcMaxLength {
			rl.tcMaxLength = l
		}
	}

	if rl.tcMaxLength > rl.MaxTabItemLength && rl.MaxTabItemLength > 0 && rl.MaxTabItemLength > rl.MinTabItemLength {
		rl.tcMaxLength = rl.MaxTabItemLength
	}
	if rl.tcMaxLength == 0 {
		rl.tcMaxLength = 20
	}

	rl.tcPosX = 1
	rl.tcPosY = 1
	rl.tcMaxX = rl.termWidth() / (rl.tcMaxLength + 2)
	rl.tcOffset = 0

	// avoid a divide by zero error
	if rl.tcMaxX < 1 {
		rl.tcMaxX = 1
	}

	rl.tcMaxY = rl.MaxTabCompleterRows

	// pre-cache
	max := rl.tcMaxX * rl.tcMaxY
	if max > len(rl.tcSuggestions) {
		max = len(rl.tcSuggestions)
	}
	subset := rl.tcSuggestions[:max]

	if rl.tcr == nil || rl.tcr.HintCache == nil {
		return
	}

	go rl.tabHintCache(subset)
}

func (rl *Instance) tabHintCache(subset []string) {
	hints := rl.tcr.HintCache(rl.tcPrefix, subset)
	if len(hints) != len(subset) {
		return
	}

	rl.tabMutex.Lock()
	for i := range subset {
		rl.tcDescriptions[subset[i]] = hints[i]
	}
	rl.tabMutex.Unlock()

}

func (rl *Instance) moveTabGridHighlight(x, y int) {
	rl.tabMutex.Lock()
	defer rl.tabMutex.Unlock()

	var suggestions *suggestionsT
	if rl.modeTabFind {
		suggestions = newSuggestionsT(rl, rl.tfSuggestions)
	} else {
		suggestions = newSuggestionsT(rl, rl.tcSuggestions)
	}

	rl.tcPosX += x
	rl.tcPosY += y

	if rl.tcPosX < 1 {
		rl.tcPosX = rl.tcMaxX
		rl.tcPosY--
	}

	if rl.tcPosX > rl.tcMaxX {
		rl.tcPosX = 1
		rl.tcPosY++
	}

	if rl.tcPosY < 1 {
		rl.tcPosY = rl.tcUsedY
	}

	if rl.tcPosY > rl.tcUsedY {
		rl.tcPosY = 1
	}

	if rl.tcPosY == rl.tcUsedY && (rl.tcMaxX*(rl.tcPosY-1))+rl.tcPosX > suggestions.Len() {
		if x < 0 {
			rl.tcPosX = suggestions.Len() - (rl.tcMaxX * (rl.tcPosY - 1))
		}

		if x > 0 {
			rl.tcPosX = 1
			rl.tcPosY = 1
		}

		if y < 0 {
			rl.tcPosY--
		}

		if y > 0 {
			rl.tcPosY = 1
		}
	}
}

func (rl *Instance) writeTabGridStr() string {
	rl.tabMutex.Lock()
	defer rl.tabMutex.Unlock()

	var suggestions *suggestionsT
	if rl.modeTabFind {
		suggestions = newSuggestionsT(rl, rl.tfSuggestions)
	} else {
		suggestions = newSuggestionsT(rl, rl.tcSuggestions)
	}

	iCellWidth := (rl.termWidth() / rl.tcMaxX) - 2
	cellWidth := strconv.Itoa(iCellWidth)

	x := 0
	y := 1
	rl.previewItem = ""
	var output string

	for i := 0; i < suggestions.Len(); i++ {
		x++
		if x > rl.tcMaxX {
			x = 1
			y++
			if y > rl.tcMaxY {
				y--
				break
			} else {
				output += "\r\n"
			}
		}

		if x == rl.tcPosX && y == rl.tcPosY {
			output += seqBgWhite + seqFgBlack
			rl.previewItem = suggestions.ItemValue(i)
		}

		value := suggestions.ItemValue(i)
		caption := cropCaption(value, rl.tcMaxLength, iCellWidth)
		if caption != value {
			rl.tcDescriptions[suggestions.ItemLookupValue(i)] = value
		}

		output += fmt.Sprintf(" %-"+cellWidth+"s %s", caption, seqReset)
	}

	rl.tcUsedY = y

	return output
}

func cropCaption(caption string, tcMaxLength int, iCellWidth int) string {
	switch {
	case iCellWidth == 0:
		// this condition shouldn't ever happen but lets cover it just in case
		return ""

	case runewidth.StringWidth(caption) != len(caption):
		// string length != rune width. So lets not do anything too clever
		//return runewidth.Truncate(caption, iCellWidth, "…")
		return runeWidthTruncate(caption, iCellWidth)

	case len(caption) < tcMaxLength,
		len(caption) < 5,
		len(caption) <= iCellWidth:
		return caption

	case len(caption)-iCellWidth+6 < 1:
		// truncate the end
		return caption[:iCellWidth-1] + "…"

	case len(caption) > 5+len(caption)-iCellWidth+6:
		// truncate long lines in the middle
		return caption[:5] + "…" + caption[len(caption)-iCellWidth+6:]

	default:
		// edge case reached. lets truncate the most conservative way we can,
		// just in case
		return runewidth.Truncate(caption, iCellWidth, "…")
	}
}
