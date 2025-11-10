package find

import (
	"strings"
)

type fuzzyFindT struct {
	mode   int
	tokens []string
}

const (
	ffMatchAll  = 0
	ffMatchSome = iota + 1
	ffMatchNone
)

func (ff *fuzzyFindT) MatchString(item string) bool {
	switch ff.mode {

	case ffMatchSome:
		return ff.matchSome(item)

	case ffMatchNone:
		return ff.matchNone(item)

	default:
		return ff.matchAll(item)
	}
}

func (ff *fuzzyFindT) Description() string {
	return "partial word"
}

func (ff *fuzzyFindT) matchAll(item string) bool {
	if len(ff.tokens) == 0 {
		return true
	}

	for i := range ff.tokens {
		if !strings.Contains(strings.ToLower(item), ff.tokens[i]) {
			return false
		}
	}

	return true
}

func (ff *fuzzyFindT) matchSome(item string) bool {
	if len(ff.tokens) == 0 {
		return true
	}

	for i := range ff.tokens {
		if strings.Contains(strings.ToLower(item), ff.tokens[i]) {
			return true
		}
	}

	return false
}

func (ff *fuzzyFindT) matchNone(item string) bool {
	if len(ff.tokens) == 0 {
		return false
	}

	for i := range ff.tokens {
		if strings.Contains(strings.ToLower(item), ff.tokens[i]) {
			return false
		}
	}

	return true
}
