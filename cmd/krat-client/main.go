package main

import (
	"fmt"
	"os"

	"github.com/humbertovnavarro/krat/pkg/client"
	"github.com/sirupsen/logrus"
)

func main() {
	if os.Getenv("DEBUG") != "" {
		err := client.StartDebug()
		if err != nil {
			logrus.Error(err)
		}
		return
	}
	err := client.Start()
	if err != nil {
		fmt.Println(err)
	}
}
