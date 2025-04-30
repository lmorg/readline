package readline

import (
	"context"
	"sync/atomic"

	"github.com/lmorg/murex/utils/lists"
)

func delayedSyntaxTimer(rl *Instance, i int32) {
	if rl.PasswordMask != 0 || rl.DelayedSyntaxWorker == nil {
		return
	}

	if rl.cacheSyntax.Get(rl.line.Runes()) != "" {
		return
	}

	if rl.line.CellLen()+rl.promptLen > rl.termWidth {
		// line wraps, which is hard to do with random ANSI escape sequences
		// so better we don't bother trying.
		return
	}

	newLine := rl.DelayedSyntaxWorker(rl.line.Runes())
	var sLine string

	if rl.SyntaxHighlighter != nil {
		sLine = rl.SyntaxHighlighter(newLine)
	} else {
		sLine = string(newLine)
	}
	rl.cacheSyntax.Append(rl.line.Runes(), sLine)

	if atomic.LoadInt32(&rl.delayedSyntaxCount) != i {
		return
	}

	output := rl.moveCursorToStartStr()
	output += sLine
	output += rl.moveCursorFromEndToLinePosStr()
	rl.print(output)
}

// DelayedTabContext is a custom context interface for async updates to the tab completions
type DelayedTabContext struct {
	rl      *Instance
	Context context.Context
	cancel  context.CancelFunc
}

// AppendSuggestions updates the tab completions with additional suggestions asynchronously
func (dtc *DelayedTabContext) AppendSuggestions(suggestions []string) {
	if dtc == nil || dtc.rl == nil {
		return
	}

	if !dtc.rl.modeTabCompletion {
		return
	}

	max := dtc.rl.MaxTabCompleterRows * 20

	if len(dtc.rl.tcSuggestions) == 0 {
		dtc.rl.ForceHintTextUpdate(" ")
	}

	dtc.rl.tabMutex.Lock()

	if dtc.rl.tcDescriptions == nil {
		dtc.rl.tcDescriptions = make(map[string]string)
	}

	for i := range suggestions {
		select {
		case <-dtc.Context.Done():
			dtc.rl.tabMutex.Unlock()
			return

		default:
			if dtc.rl.tcDescriptions[suggestions[i]] != "" ||
				(len(dtc.rl.tcSuggestions) < max && lists.Match(dtc.rl.tcSuggestions, suggestions[i])) {
				// dedup
				continue
			}
			dtc.rl.tcDescriptions[suggestions[i]] = dtc.rl.tcPrefix + suggestions[i]
			dtc.rl.tcSuggestions = append(dtc.rl.tcSuggestions, suggestions[i])
		}
	}

	dtc.rl.tabMutex.Unlock()

	output := dtc.rl.clearHelpersStr()
	//dtc.rl.ForceHintTextUpdate(" ")
	output += dtc.rl.renderHelpersStr()
	dtc.rl.print(output)
}

// AppendDescriptions updates the tab completions with additional suggestions + descriptions asynchronously
func (dtc *DelayedTabContext) AppendDescriptions(suggestions map[string]string) {
	if dtc.rl == nil {
		// This might legitimately happen with some tests
		return
	}

	if !dtc.rl.modeTabCompletion {
		return
	}

	max := dtc.rl.MaxTabCompleterRows * 20

	if len(dtc.rl.tcSuggestions) == 0 {
		dtc.rl.ForceHintTextUpdate(" ")
	}

	dtc.rl.tabMutex.Lock()

	for k := range suggestions {
		select {
		case <-dtc.Context.Done():
			dtc.rl.tabMutex.Unlock()
			return

		default:
			if dtc.rl.tcDescriptions[k] != "" ||
				(len(dtc.rl.tcSuggestions) < max && lists.Match(dtc.rl.tcSuggestions, k)) {
				// dedup
				continue
			}
			dtc.rl.tcDescriptions[k] = suggestions[k]
			dtc.rl.tcSuggestions = append(dtc.rl.tcSuggestions, k)
		}
	}

	dtc.rl.tabMutex.Unlock()

	output := dtc.rl.clearHelpersStr()
	//dtc.rl.ForceHintTextUpdate(" ")
	output += dtc.rl.renderHelpersStr()
	dtc.rl.print(output)
}

func delayedPreviewTimer(rl *Instance, fn PreviewFuncT, size *PreviewSizeT, item string) {
	var ctx context.Context

	callback := func(lines []string, pos int, err error) {
		if pos == -1 {
			if rl.previewCache != nil && rl.previewCache.pos < len(lines) {
				pos = rl.previewCache.pos
			} else {
				pos = 0
			}
		}

		select {
		case <-ctx.Done():
			return
		default:
			// continue
		}

		if err != nil {
			rl.ForceHintTextUpdate(err.Error())
			return
		}

		rl.previewCache = &previewCacheT{
			item:  item,
			pos:   pos,
			len:   size.Height,
			lines: lines,
			size:  size,
		}

		output, err := rl.previewDrawStr(lines[pos:], size)

		if err != nil {
			rl.previewCache = nil
			rl.print(output)
			return
		}

		rl.print(output)
	}

	ctx, rl.previewCancel = context.WithCancel(context.Background())
	fn(ctx, rl.line.Runes(), item, rl.PreviewImages, size, callback)
}
