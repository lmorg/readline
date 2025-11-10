package find

import "strings"

type FindT interface {
	MatchString(string) bool
	Description() string
}

func New(pattern string) (FindT, error) {
	pattern = strings.ToLower(pattern)
	ff := new(fuzzyFindT)

	ff.tokens = strings.Split(pattern, " ")

	for {
		if len(ff.tokens) == 0 {
			return ff, nil
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
		return newRegexFind(strings.Join(ff.tokens[1:], " "))

	case "g":
		pattern = strings.Join(ff.tokens[1:], " ")
		return newGlobFind(pattern)

	default:
		if strings.Contains(pattern, "*") {
			return newGlobFind(pattern)
		}
	}

	return ff, nil
}
