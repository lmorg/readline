package readline

import "testing"

// TestStrLenStripsAnsi covers every escape family strLen needs to ignore
// when computing on-screen width. Each case carries a short description
// explaining what the input represents in the wild.
func TestStrLenStripsAnsi(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		// Baselines: bytes that should be counted as-is.
		{"plain ascii", "hello", 5},
		{"empty", "", 0},
		{"unicode powerline glyph", "\ue0b0", 1},

		// SGR (colors) — what strLen handled before this fix.
		{"sgr reset", "\x1b[0mhi", 2},
		{"sgr 24-bit fg", "\x1b[38;2;255;71;156mword", 4},
		{"sgr 24-bit fg+bg", "\x1b[48;2;0;0;0m\x1b[38;2;255;255;255mok\x1b[0m", 2},

		// OSC 8 hyperlinks (oh-my-posh emits these around clickable
		// folder/git segments). Visible content here is " src ".
		{
			name: "osc 8 hyperlink ST terminator",
			in:   "\x1b]8;;https://example.com\x1b\\ src \x1b]8;;\x1b\\",
			want: 5,
		},
		{
			name: "osc 8 hyperlink BEL terminator",
			in:   "\x1b]8;;https://example.com\x07 src \x1b]8;;\x07",
			want: 5,
		},

		// OSC 0/2 window title — emitted by omp's ConsoleTitleTemplate.
		{"osc 0 title", "\x1b]0;murex in oh-my-posh\x07prompt", 6},

		// OSC 7 CWD report — emitted by many shells/prompts for terminal
		// "new tab in same dir" integrations.
		{"osc 7 cwd", "\x1b]7;file:///home/u\x1b\\prompt", 6},

		// OSC 1337 (iTerm proprietary) — CurrentDir / RemoteHost markers.
		{"osc 1337 iterm", "\x1b]1337;CurrentDir=/tmp\x07$ ", 2},

		// DEC private cursor save/restore — used by some transient
		// prompts to redraw without scrolling.
		{"dec save+restore", "\x1b7$ \x1b8", 2},

		// Cursor positioning CSI sequences other than SGR — should also
		// occupy zero visible cells.
		{"csi cursor up", "\x1b[2A>", 1},
		{"csi clear line", "\x1b[2K> ", 2},

		// Realistic compound: SGR color + OSC 8 hyperlink + glyph + text.
		// Visible content: " repo  main"
		{
			name: "compound omp-style segment",
			in:   "\x1b[38;2;195;134;241m\x1b[48;2;255;71;156m" +
				"\x1b]8;;https://github.com/user/repo\x1b\\ repo \x1b]8;;\x1b\\" +
				"\x1b[0m\ue0b0\x1b[38;2;255;251;56m  main",
			want: 13,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := strLen(tt.in)
			if got != tt.want {
				t.Errorf("strLen(%q) = %d, want %d", tt.in, got, tt.want)
			}
		})
	}
}
