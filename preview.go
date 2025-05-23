package readline

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

type previewModeT int

const (
	previewModeClosed       previewModeT = 0
	previewModeOpen         previewModeT = 1
	previewModeAutocomplete previewModeT = 2
)

type previewRefT int

const (
	previewRefDefault previewRefT = 0
	previewRefLine    previewRefT = 1
)

const (
	previewHeadingHeight = 3
	previewPromptHSpace  = 3
)

const (
	boxTL = "┏"
	boxTR = "┓"
	boxBL = "┗"
	boxBR = "┛"
	boxH  = "━"
	boxHN = "─" // narrow
	boxHD = "╶" // dashed
	boxV  = "┃"
	boxVL = "┠"
	boxVR = "┨"
)

const (
	headingTL = "╔"
	headingTR = "╗"
	headingBL = "╚"
	headingBR = "╝"
	headingH  = "═"
	headingV  = "║"
	headingVL = "╟"
	headingVR = "╢"
)

const (
	glyphScrollBar = "█"
)

// previewPos should be a percentage represented as a decimal value (eg 0.5 == 50%)
func getScrollBarSize(previewHeight int, previewPos float64) int {
	size := int((float64(previewHeight) + 2) * previewPos)
	/*if previewPos < 1 && size >= previewHeight {
		size--
	}*/
	return size
}

func getPreviewPos(rl *Instance) float64 {
	if rl.previewCache == nil {
		return 0
	}

	return (float64(rl.previewCache.pos) + float64(rl.previewCache.size.Height) + 2) / float64(len(rl.previewCache.lines))
}

func getPreviewWidth(width int) (preview, forward int) {
	preview = width - 3

	forward = width - preview
	forward -= 2
	return
}

type PreviewSizeT struct {
	Height  int
	Width   int
	Forward int
}

type previewCacheT struct {
	item  string
	pos   int
	len   int
	lines []string
	size  *PreviewSizeT
}

func (rl *Instance) getPreviewXY() (*PreviewSizeT, error) {
	if rl.isNoTty {
		return nil, errors.New("not supported in no TTY mode")
	}

	width, height, err := GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return nil, err
	}

	if height == 0 {
		height = 25
	}

	if width == 0 {
		width = 80
	}

	rl.previewAutocompleteHeight(height)

	preview, forward := getPreviewWidth(width)
	size := &PreviewSizeT{
		Height:  height - rl.MaxTabCompleterRows - 10, // hintText, multi-line prompts, etc
		Width:   preview,
		Forward: forward,
	}

	return size, nil
}

func (rl *Instance) previewAutocompleteHeight(height int) {
	switch {
	case height < 20:
		rl.MaxTabCompleterRows = noLargerThan(rl.MaxTabCompleterRows, 2)
	case height < 25:
		rl.MaxTabCompleterRows = noLargerThan(rl.MaxTabCompleterRows, 3)
	case height < 28:
		rl.MaxTabCompleterRows = noLargerThan(rl.MaxTabCompleterRows, 4)
	case height < 31:
		rl.MaxTabCompleterRows = noLargerThan(rl.MaxTabCompleterRows, 5)
	case height < 34:
		rl.MaxTabCompleterRows = noLargerThan(rl.MaxTabCompleterRows, 6)
	case height < 37:
		rl.MaxTabCompleterRows = noLargerThan(rl.MaxTabCompleterRows, 7)
	case height < 40:
		rl.MaxTabCompleterRows = noLargerThan(rl.MaxTabCompleterRows, 8)
	}
}

func noLargerThan(src, max int) int {
	if src > max {
		return max
	}
	return src
}

func (rl *Instance) writePreviewStr() string {
	if rl.previewMode == previewModeClosed {
		rl.previewCache = nil
		return ""
	}

	if rl.previewCancel != nil {
		rl.previewCancel()
	}

	var fn PreviewFuncT
	if rl.previewRef == previewRefLine {
		fn = rl.PreviewLine
	} else {
		if rl.tcr == nil {
			rl.previewCache = nil
			return ""
		}
		fn = rl.tcr.Preview
	}

	if fn == nil {
		rl.previewCache = nil
		return ""
	}

	size, err := rl.getPreviewXY()
	if err != nil || size.Height < 4 || size.Width < 10 {
		rl.previewCache = nil
		return previewTerminalTooSmall
	}

	item := rl.previewItem
	item = strings.ReplaceAll(item, "\\", "")
	item = strings.TrimSpace(item)

	go delayedPreviewTimer(rl, fn, size, item)

	return ""
}

var previewTerminalTooSmall = fmt.Sprintf("%s%sTerminal too small to display preview%s", curPosSave, curHome, curPosRestore)

const (
	curHome       = "\x1b[H"
	curPosSave    = "\x1b[s"
	curPosRestore = "\x1b[u"
)

func (rl *Instance) previewDrawStr(preview []string, size *PreviewSizeT) (string, error) {
	var (
		output       string
		scrollBar    = glyphScrollBar
		scrollHeight = getScrollBarSize(size.Height-1, getPreviewPos(rl))
	)

	pf := fmt.Sprintf("%s%%-%ds%s\r\n", boxV, size.Width, scrollBar)
	pj := fmt.Sprintf("%s%%-%ds%s\r\n", boxVL, size.Width, scrollBar)

	output += curHome

	output += fmt.Sprintf(cursorForwf, size.Forward)
	hr := strings.Repeat(headingH, size.Width)
	output += headingTL + hr + headingTR + "\r\n "
	output += headingV + rl.previewTitleStr(size.Width) + headingV + "\r\n "
	output += headingBL + hr + headingBR + "\r\n "

	hr = strings.Repeat(boxH, size.Width)
	output += boxTL + hr + boxTR + "\r\n"

	for i := 0; i <= size.Height; i++ {
		if i == scrollHeight {
			scrollBar = boxV
			pf = fmt.Sprintf("%s%%-%ds%s\r\n", boxV, size.Width, scrollBar)
			pj = fmt.Sprintf("%s%%-%ds%s\r\n", boxVL, size.Width, boxVR)
		}

		output += fmt.Sprintf(cursorForwf, size.Forward)

		if i >= len(preview) {
			blank := strings.Repeat(" ", size.Width)
			output += boxV + blank + scrollBar + "\r\n"
			continue
		}

		if strings.HasPrefix(preview[i], boxHN) || strings.HasPrefix(preview[i], boxHD) {
			output += fmt.Sprintf(pj, preview[i])
		} else {
			output += fmt.Sprintf(pf, preview[i])
		}
	}

	output += fmt.Sprintf(cursorForwf, size.Forward)
	output += boxBL + hr + boxBR + "\r\n"

	output += rl.previewMoveToPromptStr(size)
	return output, nil
}

func (rl *Instance) previewTitleStr(width int) string {
	var title string

	if rl.previewRef == previewRefDefault {
		title = " Autocomplete Preview" + title
	} else {
		title = " Command Line Preview" + title
	}
	title += "    |    [F1] to exit    |    [ENTER] to commit"

	l := len(title) + 1
	switch {
	case l > width:
		return title[:width-2] + "… "
	case l == width:
		return title + " "
	default:
		return title + strings.Repeat(" ", width-l+1)
	}
}

func (rl *Instance) previewMoveToPromptStr(size *PreviewSizeT) string {
	output := curHome
	output += moveCursorDownStr(size.Height + previewPromptHSpace + previewHeadingHeight)
	output += rl.moveCursorFromStartToLinePosStr()
	return output
}

func (rl *Instance) previewPreviousSectionStr() string {
	if rl.previewCache == nil || rl.previewCache.pos == 0 {
		return ""
	}

	for rl.previewCache.pos -= 2; rl.previewCache.pos > 0; rl.previewCache.pos-- {
		if strings.HasPrefix(rl.previewCache.lines[rl.previewCache.pos], boxHN) {
			if rl.previewCache.pos < len(rl.previewCache.lines)-1 {
				rl.previewCache.pos++
			}
			break
		}
	}

	if rl.previewCache.pos > len(rl.previewCache.lines)-rl.previewCache.len-1 {
		rl.previewCache.pos = len(rl.previewCache.lines) - rl.previewCache.len - 1
	}
	if rl.previewCache.pos < 0 {
		rl.previewCache.pos = 0
	}

	output, _ := rl.previewDrawStr(rl.previewCache.lines[rl.previewCache.pos:], rl.previewCache.size)
	return output
}

func (rl *Instance) previewNextSectionStr() string {
	if rl.previewCache == nil {
		return ""
	}

	for ; rl.previewCache.pos < len(rl.previewCache.lines)-rl.previewCache.len; rl.previewCache.pos++ {
		if strings.HasPrefix(rl.previewCache.lines[rl.previewCache.pos], boxHN) {
			if rl.previewCache.pos < len(rl.previewCache.lines)-1 {
				rl.previewCache.pos++
			}
			break
		}
	}

	if rl.previewCache.pos > len(rl.previewCache.lines)-rl.previewCache.len-1 {
		rl.previewCache.pos = len(rl.previewCache.lines) - rl.previewCache.len - 1
	}
	if rl.previewCache.pos < 0 {
		rl.previewCache.pos = 0
	}

	output, _ := rl.previewDrawStr(rl.previewCache.lines[rl.previewCache.pos:], rl.previewCache.size)
	return output
}

func (rl *Instance) previewPageUpStr() string {
	if rl.previewCache == nil {
		return ""
	}

	rl.previewCache.pos -= rl.previewCache.len
	if rl.previewCache.pos < 0 {
		rl.previewCache.pos = 0
	}

	output, _ := rl.previewDrawStr(rl.previewCache.lines[rl.previewCache.pos:], rl.previewCache.size)
	return output
}

func (rl *Instance) previewPageDownStr() string {
	if rl.previewCache == nil {
		return ""
	}

	rl.previewCache.pos += rl.previewCache.len
	if rl.previewCache.pos > len(rl.previewCache.lines)-rl.previewCache.len-1 {
		rl.previewCache.pos = len(rl.previewCache.lines) - rl.previewCache.len - 1
		if rl.previewCache.pos < 0 {
			rl.previewCache.pos = 0
		}
	}

	output, _ := rl.previewDrawStr(rl.previewCache.lines[rl.previewCache.pos:], rl.previewCache.size)
	return output
}

func (rl *Instance) clearPreviewStr() string {
	var output string

	if rl.previewCancel != nil {
		rl.previewCancel()
	}

	if rl.PreviewInit != nil {
		rl.PreviewInit()
	}

	if rl.previewMode > previewModeClosed {
		output = seqRestoreBuffer + curPosRestore
		output += rl.echoStr()
		rl.previewMode = previewModeClosed
		rl.previewRef = previewRefDefault
	}

	return output
}
