package cli

import (
	"github.com/CavemanJay/gogurt/config"

	"github.com/urfave/cli/v2"
)

var configuration = config.Config{}

func InitApp(version string) *cli.App {
	app := cli.NewApp()
	app.Name = "Gogurt"
	app.Usage = "Handles functions I commonly need to do"
	app.UseShortOptionHandling = true
	app.EnableBashCompletion = true
	app.Version = version
	app.Flags = appFlags
	app.Commands = appCommands

	return app
}
