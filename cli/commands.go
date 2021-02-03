package cli

import "github.com/urfave/cli/v2"

var appCommands = []*cli.Command{
	{
		Name:    "listen",
		Aliases: []string{"l"},
		Usage:   "Listen on the specified port for clients",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:        "port",
				Aliases:     []string{"p"},
				Usage:       "The `PORT` to listen on",
				Value:       8081,
				Destination: &configuration.Port,
			},
		},
		Action: listen,
	},
	{
		Name:    "dial",
		Aliases: []string{"d"},
		Usage:   "Connects to an existing server instance",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:        "port",
				Aliases:     []string{"p"},
				Usage:       "The `PORT` to connect to",
				Value:       8081,
				Destination: &configuration.Port,
			},
			&cli.StringFlag{
				Name:        "host",
				Usage:       "The `HOST` to connect to",
				Destination: &configuration.ClientOptions.Host,
				Value:       "localhost",
			},
		},
		Action: dial,
	},
	{
		Name:    "archive",
		Aliases: []string{"a"},
		Action:  archiveAction,
	},
	{
		Name:    "sync-music",
		Aliases: []string{"m"},
		Action:  syncMusic,
	},
}
