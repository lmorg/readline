package readline

import (
	"testing"
)

// TestViDeleteByAdjustLogicNoPanic just tests that viDeleteByAdjust() doesn't cause
// a panic:
// https://github.com/lmorg/murex/issues/341

type TestViDeleteByAdjustT struct {
	Line   string
	Pos    int
	Adjust int
}

func TestViDeleteByAdjustLogicNoPanic(t *testing.T) {
	tests := []TestViDeleteByAdjustT{
		{
			Line:   "The quick brown fox",
			Pos:    0,
			Adjust: -1,
		},
		{
			Line:   "The quick brown fox",
			Pos:    1,
			Adjust: -1,
		},
		{
			Line:   "The quick brown fox",
			Pos:    1,
			Adjust: -2,
		},
		{
			Line:   "The quick brown fox",
			Pos:    2,
			Adjust: -2,
		},
		{
			Line:   "The quick brown fox",
			Pos:    5,
			Adjust: -1,
		},
		{
			Line:   "The quick brown fox",
			Pos:    5,
			Adjust: 1,
		},
		{
			Line:   "The quick brown fox",
			Pos:    5,
			Adjust: 100,
		},
	}

	for _, test := range tests {
		rl := NewInstance()
		rl.line.Set(rl, []rune(test.Line))
		rl.line.SetRunePos(test.Pos)
		rl.viDeleteByAdjustLogic(&test.Adjust)
	}
}
