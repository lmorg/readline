//go:build !js
// +build !js

package readline

import "os"

func _print(s string) {
	os.Stdout.WriteString(s)
}

func _printErr(s string) {
	os.Stderr.WriteString(s)
}
