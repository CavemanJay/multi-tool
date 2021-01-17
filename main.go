package main

import (
	"log"
	"os"

	"github.com/JayCuevas/jays-server/cli"
)

func main() {
	app := cli.InitApp()
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
