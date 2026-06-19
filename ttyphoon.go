package readline

import (
	"fmt"
	"os"
	"strings"
)

func isTtyphoon() bool {
	return os.Getenv("MXTTY") == "true"
}

const (
	seqApc = "\x1b_"
	seqST  = "\x1b\\"
)

func apcReply(r []rune) (string, string) {
	s := strings.SplitN(string(r), ";", 3)

	if len(s) < 3 || s[0] != seqApc+"reply" {
		return "", ""
	}

	return s[1], s[2][:len(s[2])-2]
}

func (rl *Instance) apcSequence(r []rune) string {
	key, value := apcReply(r)
	switch key {
	case "content-editable":
		return rl.apcContentEditable(value)

	case "":
		fallthrough
	default:
		return ""
	}
}

func (rl *Instance) apcContentEditable(s string) string {
	s = strings.ReplaceAll(s, strings.Repeat(" ", len(rl.prompt)), "")
	rl.line.Set(rl, []rune(s))
	if rl.line.RunePos() > len(s) {
		rl.line.SetRunePos(len(s))
	}

	output := rl.moveCursorByRuneAdjustStr(len(s))
	output += rl.echoStr()
	output += rl.updateHelpersStr()
	return output
}

func contentEditable(s string) string {
	return fmt.Sprintf("%sbegin;content-editable%s%s%send;content-editable;%s%s",
		seqApc, seqST, //                                 // opening tag
		s,                                                   // content editable content
		seqApc, rxAnsiEscape.ReplaceAllString(s, ""), seqST, // closing tag
	)
}
