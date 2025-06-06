//go:build js
// +build js

package readline

import "syscall/js"

func _print(s string) {
	vtermWrite([]rune(s))
}

func _printErr(s string) {
	vtermWrite([]rune(s))
}

func vtermWrite(r []rune) {
	VTerm.Write(r)

	//divMutex.Lock()

	html := VTerm.ExportHtml()

	jsDoc := js.Global().Get("document")
	outElement := jsDoc.Call("getElementById", "term")
	outElement.Set("innerHTML", html)

	//divMutex.Unlock()
}
