# lmorg/readline

v4.0.0 marks a breaking change to the tab completion function.

Earlier versions expected multiple parameters to be returned however from
v4.0.0 onwards, a pointer to a structure is instead expected:
```
type TabCompleterReturnT struct {
	Prefix       string
	Suggestions  []string
	Descriptions map[string]string
	DisplayType  TabDisplayType
	HintCache    HintCacheFuncT
	Preview      PreviewFuncT
}
```
This allows for more configurability and without the cost of copying multiple
different pieces of data nor future breaking changes whenever additional new
features are added.

## Changes

### 4.1.0

* Murex has switched back to calling this package for `readline`, meaning this
  package will now see more regular updates and bug fixes

* bugfix: cursor wouldn't step backwards in VIM mode when cursor at end of line

* experimental support added for integrating `readline` into GUI applications

### 4.0.0

* support for wide and zero width unicode characters
  ([inherited from Murex](https://murex.rocks/changelog/v4.0.html))

* preview modes
  ([inherited from Murex](https://murex.rocks/user-guide/interactive-shell.html#preview))

* API improvements

* rewritten event system
  ([discussion](https://github.com/lmorg/murex/discussions/799))

* vastly improved buffered rendering -- this leads to few rendering glitches
  and particularly on slower machines and/or terminals

* added missing vim and emacs keybindings
  ([full list of keybindings](listhttps://murex.rocks/user-guide/terminal-keys.html))

* additional tests

* fixed glitches on Windows terminals
  ([discussion](https://github.com/lmorg/murex/issues/630))

* readline command mode
  ([discussion](https://github.com/lmorg/murex/discussions/905))

### 3.0.1

This is a bug fix release:

* Nil map panic fixed when using dtx.AppendSuggestions()

* Hint text line proper blanked (this is a fix to a regression bug introduced
  in version 3.0.0)

* Example 01 updated to reflect API changes in 3.0.0

### 3.0.0

This release brings a considerable number of new features and bug fixes
inherited from readline's use in murex (https://github.com/lmorg/murex)

* Wrapped lines finally working (where the input line is longer than the
  terminal width)

* Delayed tab completion - allows asynchronous updates to the tab completion so
  slower suggestions do not halt the user experience

* Delayed syntax timer - allows syntax highlighting to run asynchronously for
  slower parsers (eg spell checkers)

* Support for GetCursorPos ANSI escape sequence (though I don't have a terminal
  which supports this to test the code on)

* Better support for wrapped hint text lines

* Fixed bug with $EDITOR error handling in Windows and Plan 9

* Code clean up - fewer writes to the terminal

If you just use the exported API end points then your code should still work
verbatim. However if you are working against a fork or custom patch set then
considerable more work may be required to merge the changes.

### 2.1.0

Error returns from `readline` have been created as error a variable, which is
more idiomatic to Go than the err constants that existed previously. Currently
both are still available to use however I will be deprecating the the constants
in a latter release.

**Deprecated constants:**
```go
const (
	// ErrCtrlC is returned when ctrl+c is pressed
	ErrCtrlC = "Ctrl+C"

	// ErrEOF is returned when ctrl+d is pressed
	ErrEOF = "EOF"
)
```

**New error variables:**
```go
var (
	// CtrlC is returned when ctrl+c is pressed
	CtrlC = errors.New("Ctrl+C")

	// EOF is returned when ctrl+d is pressed
	// (this is actually the same value as io.EOF)
	EOF = errors.New("EOF")
)
```
