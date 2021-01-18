package main

import (
	"log"
	"os"

	"github.com/JayCuevas/gogurt/cli"
)

func main() {
	app := cli.InitApp()
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
