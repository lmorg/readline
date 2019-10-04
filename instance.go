package readline

import (
	"os"
	"sync"
)

// Instance is used to encapsulate the parameter group and run time of any given
// readline instance so that you can reuse the readline API for multiple entry
// captures without having to repeatedly unload configuration.
type Instance struct {
	mutex sync.Mutex

	// PasswordMask is what character to hide password entry behind.
	// Once enabled, set to 0 (zero) to disable the mask again.
	PasswordMask rune

	// SyntaxHighlight is a helper function to provide syntax highlighting.
	// Once enabled, set to nil to disable again.
	SyntaxHighlighter func([]rune) string

	// History is an interface for querying the readline history.
	// This is exposed as an interface to allow you the flexibility to define how
	// you want your history managed (eg file on disk, database, cloud, or even
	// no history at all). By default it uses a dummy interface that only stores
	// historic items in memory.
	History History

	// HistoryAutoWrite defines whether items automatically get written to
	// history.
	// Enabled by default. Set to false to disable.
	HistoryAutoWrite bool // = true

	// TabCompleter is a simple function that offers completion suggestions.
	// It takes the readline line ([]rune) and cursor pos. Returns a prefix
	// string, an array of suggestions and a map of definitions (optional).
	TabCompleter      func([]rune, int, DelayedTabContext) (string, []string, map[string]string, TabDisplayType)
	delayedTabContext DelayedTabContext

	// MaxTabCompletionRows is the maximum number of rows to display in the tab
	// completion grid.
	MaxTabCompleterRows int // = 4

	// SyntaxCompletion is used to autocomplete code syntax (like braces and
	// quotation marks). If you want to complete words or phrases then you might
	// be better off using the TabCompletion function.
	// SyntaxCompletion takes the line ([]rune) and cursor position, and returns
	// the new line and cursor position.
	SyntaxCompleter func([]rune, int) ([]rune, int)

	// DelayedSyntaxWorker allows for syntax highlighting happen to the line
	// after the line has been drawn.
	DelayedSyntaxWorker func([]rune) []rune
	delayedSyntaxCount  int64

	// HintText is a helper function which displays hint text the prompt.
	// HintText takes the line input from the prompt and the cursor position.
	// It returns the hint text to display.
	HintText func([]rune, int) []rune

	// HintColor any ANSI escape codes you wish to use for hint formatting. By
	// default this will just be blue.
	HintFormatting string

	// TempDirectory is the path to write temporary files when editing a line in
	// $EDITOR. This will default to os.TempDir()
	TempDirectory string

	// GetMultiLine is a callback to your host program. Since multiline support
	// is handled by the application rather than readline itself, this callback
	// is required when calling $EDITOR. However if this function is not set
	// then readline will just use the current line.
	GetMultiLine func([]rune) []rune

	// readline operating parameters
	prompt        string //  = ">>> "
	promptLen     int    //= 4
	line          []rune
	pos           int
	multiline     []byte
	multisplit    []string
	skipStdinRead bool

	// history
	lineBuf string
	histPos int

	// hint text
	hintY    int //= 0
	hintText []rune

	// tab completion
	modeTabCompletion bool
	tcPrefix          string
	tcSuggestions     []string
	tcDescriptions    map[string]string
	tcDisplayType     TabDisplayType
	tcOffset          int
	tcPosX            int
	tcPosY            int
	tcMaxX            int
	tcMaxY            int
	tcUsedY           int
	tcMaxLength       int

	// tab find
	modeTabFind   bool
	tfLine        []rune
	tfSuggestions []string
	modeAutoFind  bool // for when invoked via ^R or ^F outside of [tab]

	// vim
	modeViMode       viMode //= vimInsert
	viIteration      string
	viUndoHistory    []undoItem
	viUndoSkipAppend bool
	viYankBuffer     string

	// event
	evtKeyPress map[string]func(string, []rune, int) *EventReturn
}

// NewInstance is used to create a readline instance and initialise it with sane
// defaults.
func NewInstance() *Instance {
	rl := new(Instance)

	//GetTermWidth()

	rl.History = new(ExampleHistory)
	rl.HistoryAutoWrite = true
	rl.MaxTabCompleterRows = 4
	rl.prompt = ">>> "
	rl.promptLen = 4
	rl.HintFormatting = seqFgBlue
	rl.evtKeyPress = make(map[string]func(string, []rune, int) *EventReturn)

	rl.TempDirectory = os.TempDir()

	return rl
}
