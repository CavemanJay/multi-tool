package cli

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
	"path"

	"github.com/CavemanJay/gogurt/client"
	"github.com/CavemanJay/gogurt/config"
	"github.com/CavemanJay/gogurt/server"
	"github.com/op/go-logging"

	"github.com/urfave/cli/v2"
)

var configuration = config.Config{}

func getAppDataPath() string {
	appData, err := os.UserCacheDir()
	if err != nil {
		panic(err)
	}

	return path.Join(appData, "gogurt")
}

func initLogger(file io.Writer) {

	// format := logging.MustStringFormatter(`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03d}%{color:reset} %{message}`)
	// format := logging.MustStringFormatter(`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`)
	format := logging.MustStringFormatter(`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level} %{id:03x}%{color:reset} %{message}`)

	// format := logging.MustStringFormatter(`%{color}%{time:15:04:05.000} %{callpath} ▶ %{level:2.5s} %{id:03x}%{color:reset} %{message}`)

	stdOutBackend := logging.NewLogBackend(os.Stdin, "", 0)
	stdOut := logging.AddModuleLevel(logging.NewBackendFormatter(stdOutBackend, format))
	stdOut.SetLevel(logging.DEBUG, "gogurt")

	// If we are in a tty
	// if fileInfo, _ := os.Stdout.Stat(); (fileInfo.Mode() & os.ModeCharDevice) != 0 {
	// } else {
	// 	stdErr := logging.AddModuleLevel(logging.NewLogBackend(os.Stderr, "", 0))
	// 	stdErr.SetLevel(logging.ERROR, "ERROR")
	// 	logging.SetBackend(stdErr, logFile, stdOut)
	// }

	if file != nil {
		logFile := logging.AddModuleLevel(logging.NewLogBackend(file, "", log.Lshortfile|log.Ldate|log.Ltime))
		logFile.SetLevel(logging.INFO, "gogurt")
		logging.SetBackend(logFile, stdOut)
	} else {
		logging.SetBackend(stdOut)
	}

}

func InitApp() *cli.App {
	app := cli.NewApp()
	app.Name = "Gogurt"
	app.Usage = "Handles functions I commonly need to do"
	app.UseShortOptionHandling = true
	app.EnableBashCompletion = true
	app.Version = "0.1.0"

	u, err := user.Current()
	if err != nil {
		panic(err)
	}

	app.Flags = []cli.Flag{
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
			Value:       path.Join(u.HomeDir, "Sync"),
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
			Name:        "use-config",
			Aliases:     []string{"c"},
			Usage:       "The `FILE` to read config values from",
			Destination: &configuration.ConfigPath,
		},
	}

	app.Commands = []*cli.Command{
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
			// Usage:   "Archives",
			Flags:  []cli.Flag{},
			Action: archive,
		},
	}

	return app
}

func listen(ctx *cli.Context) error {
	cfg := &configuration

	if err := handleConfig(); err != nil {
		return err
	}

	_, err := os.Stat(cfg.SyncFolder)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("Folder \"%s\" does not exist", cfg.SyncFolder)
		}
		return err
	}

	os.Mkdir(cfg.AppDataFolder, os.ModePerm)

	logFile, err := os.OpenFile(path.Join(cfg.AppDataFolder, "gogurt.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer logFile.Close()
	initLogger(logFile)

	if cfg.Append != "" {
		cfg.SyncFolder = path.Join(cfg.SyncFolder, cfg.Append)
	}

	config.WriteConfig(path.Join(cfg.AppDataFolder, "last_run.json"), cfg)

	s := server.NewServer(cfg.SyncFolder, cfg.Recursive, cfg.Port)

	return s.Listen()
}

func dial(ctx *cli.Context) error {

	cfg := &configuration

	if err := handleConfig(); err != nil {
		return err
	}

	c := client.NewClient(cfg.SyncFolder)
	initLogger(nil)

	config.WriteConfig(path.Join(cfg.AppDataFolder, "last_run.json"), cfg)

	return c.Connect(cfg.ClientOptions.Host, cfg.Port)
}

func handleConfig() error {
	cfg := &configuration
	if cfg.UseLastRun {
		cfgFile := path.Join(getAppDataPath(), "last_run.json")
		cfg, err := config.ReadConfig(cfgFile)
		if err != nil {
			return err
		}
		configuration = *cfg
	} else if len(cfg.ConfigPath) > 0 {
		cfg, err := config.ReadConfig(cfg.ConfigPath)
		if err != nil {
			return err
		}
		configuration = *cfg
	}
	return nil
}

func archive(ctx *cli.Context) error {
	cfg := &configuration
	if err := handleConfig(); err != nil {
		return err
	}
	initLogger(nil)

	log.Printf("%#v", cfg)

	return nil
}
