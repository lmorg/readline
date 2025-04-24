package main

import (
	"fmt"
	"os"

	"github.com/lmorg/readline/v4"
	"golang.org/x/term"
)

func main() {
	rl := readline.NewInstance()

	ch := make(chan *readline.NoTtyCallbackT)

	rl.SetNoTtyCallback(ch)

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	go func() {
		for {
			p := make([]byte, 1024)
			i, err := os.Stdin.Read(p)
			if err != nil {
				panic(err)
			}

			rl.KeyPress(p[:i])
		}
	}()

	go func() {
		for {
			callback := <-ch
			fmt.Printf("\r\n>>> %s\r\n(%s)", callback.Line.String(), callback.Hint)
		}
	}()

	_, _ = rl.Readline()
}
