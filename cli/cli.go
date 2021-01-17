package cli

import (
	"os/user"
	"path"

	"github.com/JayCuevas/jays-server/server"

	"github.com/urfave/cli/v2"
)

type ServerOptions struct {
	Recursive  bool
	RootFolder string
}

type ClientOptions struct {
	Host string
}

type Config struct {
	Port          int
	ServerOptions ServerOptions
	ClientOptions ClientOptions
}

var Configuration = Config{}

func InitApp() *cli.App {
	app := cli.NewApp()
	app.Name = "Jay's Server"
	app.Usage = "Handles functions I commonly need to do"
	app.UseShortOptionHandling = true
	app.EnableBashCompletion = true

	u, err := user.Current()
	if err != nil {
		panic(err)
	}

	app.Commands = []*cli.Command{
		&cli.Command{
			Name:    "listen",
			Aliases: []string{"l"},
			Usage:   "Listen on the specified port for clients",
			Flags: []cli.Flag{
				&cli.IntFlag{
					Name:        "port",
					Aliases:     []string{"p"},
					Usage:       "The `PORT` to listen on",
					Value:       8080,
					Destination: &Configuration.Port,
				},
				&cli.BoolFlag{
					Name:        "recursive",
					Aliases:     []string{"r"},
					Usage:       "Whether or not to recursively watch the root folder",
					Destination: &Configuration.ServerOptions.Recursive,
				},
				&cli.PathFlag{
					Name:        "folder",
					Aliases:     []string{"f"},
					Usage:       "The root `FOLDER` to synchronize",
					Value:       path.Join(u.HomeDir, "Sync"),
					Destination: &Configuration.ServerOptions.RootFolder,
				},
			},
			Action: func(c *cli.Context) error {
				cfg := &Configuration
				s := server.NewServer(cfg.ServerOptions.RootFolder, cfg.ServerOptions.Recursive, cfg.Port)

				return s.Listen()
			},
		},
		&cli.Command{
			Name:    "dial",
			Aliases: []string{"d"},
			Usage:   "Connects to an existing server instance",
			Flags: []cli.Flag{
				&cli.IntFlag{
					Name:        "port",
					Aliases:     []string{"p"},
					Usage:       "The `PORT` to listen on",
					Value:       8080,
					Destination: &Configuration.Port,
				},
				&cli.StringFlag{
					Name:        "host",
					Aliases:     []string{"h"},
					Usage:       "The `HOST` to connect to",
					Destination: &Configuration.ClientOptions.Host,
				},
			},
		},
	}

	return app
}
