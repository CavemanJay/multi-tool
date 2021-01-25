package main

import (
	"os"

	"github.com/CavemanJay/gogurt/cli"
	"github.com/op/go-logging"
)

var (
	version string
)

func main() {
	app := cli.InitApp(version)
	err := app.Run(os.Args)
	if err != nil {
		log := logging.MustGetLogger("gogurt")
		log.Fatal(err)
	}
}
