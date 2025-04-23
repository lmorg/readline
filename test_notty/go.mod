module github.com/lmorg/readline/test_notty

go 1.24.1

require (
	github.com/lmorg/readline/v4 v4.0.1
	golang.org/x/term v0.31.0
)

require (
	github.com/lmorg/murex v0.0.0-20250115225944-b4c429617fd4 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	golang.org/x/sys v0.32.0 // indirect
)

replace github.com/lmorg/readline/v4 => ../../readline
