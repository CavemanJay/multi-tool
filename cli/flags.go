package cli

import (
	"os/user"
	"path"

	"github.com/urfave/cli/v2"
)

var (
	usr      *user.User
	appFlags []cli.Flag
)

func init() {
	var err error
	usr, err = user.Current()
	if err != nil {
		panic(err)
	}

	appFlags = []cli.Flag{
		&cli.PathFlag{
			Name:        "appdata",
			Usage:       "The `PATH` to the folder where app data is stored",
			Value:       getAppDataPath(),
			Destination: &configuration.AppDataFolder,
		},
		&cli.PathFlag{
			Name:        "folder",
			Aliases:     []string{"f"},
			Usage:       "The root `FOLDER` to synchronize",
			Value:       path.Join(usr.HomeDir, "Sync"),
			Destination: &configuration.SyncFolder,
		},
		&cli.StringFlag{
			Name:        "append",
			Aliases:     []string{"a"},
			Usage:       "Appends `PATH` to the root folder",
			Destination: &configuration.Append,
		},
		&cli.BoolFlag{
			Name:        "recursive",
			Aliases:     []string{"r"},
			Usage:       "Whether or not to recursively watch the root folder",
			Destination: &configuration.Recursive,
			Value:       true,
		},
		&cli.BoolFlag{
			Name:        "use-last-run",
			Aliases:     []string{"u"},
			Usage:       "Use the options specified in the last run",
			Destination: &configuration.UseLastRun,
		},
		&cli.PathFlag{
			Name:        "config",
			Aliases:     []string{"c"},
			Usage:       "The `FILE` to read config values from",
			Destination: &configuration.ConfigPath,
		},
	}
}
