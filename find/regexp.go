package find

import (
	"regexp"
)

type regex struct {
	rx *regexp.Regexp
}

func newRegexFind(pattern string) (*regex, error) {
	rx, err := regexp.Compile("(?i)" + pattern)
	if err != nil {
		return nil, err
	}
	return &regex{rx: rx}, nil
}

func (rx *regex) MatchString(item string) bool {
	return rx.rx.MatchString(item)
}

func (rx *regex) Description() string {
	return "regexp"
}
