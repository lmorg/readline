package find

import (
	"path"
	"strings"
)

type glob struct{ pattern string }

func newGlobFind(pattern string) (*glob, error) {
	gf := new(glob)
	gf.pattern = pattern
	return gf, nil
}

func (g *glob) MatchString(item string) bool {
	found, _ := path.Match(g.pattern, strings.ToLower(item))
	return found
}

func (g *glob) Description() string {
	return "globbing"
}
