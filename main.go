package main

import (
	"fmt"
	"os"

	"github.com/CavemanJay/gogurt/cli"
)

var (
	version string
)

func main() {
	app := cli.InitApp(version)
	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
