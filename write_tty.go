//go:build !js && !readline_notty
// +build !js,!readline_notty

package readline

import "os"

func print(s string) {
	os.Stdout.WriteString(s)
}

func printErr(s string) {
	os.Stderr.WriteString(s)
}
