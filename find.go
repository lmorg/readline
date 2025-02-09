package readline

import (
	"path"
	"regexp"
	"strings"
)

var (
	rFindSearchPart = []rune("partial word match: ")
	rFindCancelPart = []rune("Cancelled partial word match")

	rFindSearchGlob = []rune("globbing match: ")
	rFindCancelGlob = []rune("Cancelled globbing match")

	rFindSearchRegex = []rune("regexp match: ")
	rFindCancelRegex = []rune("Cancelled regexp match")
)

type findT interface {
	MatchString(string) bool
}

type fuzzyFindT struct {
	mode   int
	tokens []string
}

const (
	ffMatchAll  = 0
	ffMatchSome = iota + 1
	ffMatchNone
	ffMatchRegexp
	ffMatchGlob
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

type globFindT struct{ pattern string }

func (gf *globFindT) MatchString(item string) bool {
	found, _ := path.Match(gf.pattern, item)
	return found
}

func newGlobFind(pattern string) (*globFindT, error) {
	gf := new(globFindT)
	gf.pattern = pattern
	return gf, nil
}

func newFuzzyFind(pattern string) (findT, []rune, []rune, error) {
	pattern = strings.ToLower(pattern)
	ff := new(fuzzyFindT)

	ff.tokens = strings.Split(pattern, " ")

	for {
		if len(ff.tokens) == 0 {
			return ff, rFindSearchPart, rFindCancelPart, nil
		}

		if ff.tokens[len(ff.tokens)-1] == "" {
			ff.tokens = ff.tokens[:len(ff.tokens)-1]
		} else {
			break
		}
	}

	switch ff.tokens[0] {
	case "or":
		ff.mode = ffMatchSome
		ff.tokens = ff.tokens[1:]

	case "!":
		ff.mode = ffMatchNone
		ff.tokens = ff.tokens[1:]

	case "rx":
		ff.mode = ffMatchRegexp
		pattern = strings.Join(ff.tokens[1:], " ")
		find, err := regexp.Compile("(?i)" + pattern)
		return find, rFindSearchRegex, rFindCancelRegex, err

	case "g":
		ff.mode = ffMatchGlob
		pattern = strings.Join(ff.tokens[1:], " ")
		find, err := newGlobFind(pattern)
		return find, rFindSearchGlob, rFindCancelGlob, err

	default:
		if strings.Contains(pattern, "*") {
			ff.mode = ffMatchGlob
			find, err := newGlobFind(pattern)
			return find, rFindSearchGlob, rFindCancelGlob, err
		}
	}

	return ff, rFindSearchPart, rFindCancelPart, nil
}
