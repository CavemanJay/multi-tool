package main

import (
	"os"

	"github.com/JayCuevas/gogurt/cli"
	"github.com/op/go-logging"
)

func main() {
	app := cli.InitApp()
	err := app.Run(os.Args)
	if err != nil {
		log := logging.MustGetLogger("gogurt")
		log.Fatal(err)
	}
}
