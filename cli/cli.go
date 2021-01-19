package cli

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
	"path"

	"github.com/JayCuevas/gogurt/client"
	"github.com/JayCuevas/gogurt/server"
	"github.com/op/go-logging"

	"github.com/urfave/cli/v2"
)

type ServerOptions struct {
	Recursive  bool
	RootFolder string
	Append     string
}

type ClientOptions struct {
	Host string
}

type Config struct {
	Port          int
	ServerOptions ServerOptions
	ClientOptions ClientOptions
	AppDataFolder string
}

var Configuration = Config{}

func getAppDataPath() string {
	appData, err := os.UserCacheDir()
	if err != nil {
		panic(err)
	}

	return path.Join(appData, "gogurt")
}

func initLogger(file io.Writer) {

	logFile := logging.AddModuleLevel(logging.NewLogBackend(file, "", log.Lshortfile|log.Ldate|log.Ltime))
	logFile.SetLevel(logging.INFO, "gogurt")

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

	logging.SetBackend(logFile, stdOut)
}

func InitApp() *cli.App {
	app := cli.NewApp()
	app.Name = "Gogurt"
	app.Usage = "Handles functions I commonly need to do"
	app.UseShortOptionHandling = true
	app.EnableBashCompletion = true
	app.Version = "v0.1.0"

	u, err := user.Current()
	if err != nil {
		panic(err)
	}

	app.Flags = []cli.Flag{
		&cli.PathFlag{
			Name:        "appdata",
			Usage:       "The `PATH` to the folder where app data is stored",
			Value:       getAppDataPath(),
			Destination: &Configuration.AppDataFolder,
		},
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
					Value:       8081,
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
				&cli.PathFlag{
					Name:        "append",
					Aliases:     []string{"a"},
					Usage:       "Appends `PATH` to the root folder",
					Destination: &Configuration.ServerOptions.Append,
				},
			},
			Action: func(c *cli.Context) error {
				cfg := &Configuration
				_, err := os.Stat(cfg.ServerOptions.RootFolder)
				if err != nil {
					if os.IsNotExist(err) {
						return fmt.Errorf("Folder \"%s\" does not exist", cfg.ServerOptions.RootFolder)
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
					Usage:       "The `PORT` to connect to",
					Value:       8081,
					Destination: &Configuration.Port,
				},
				&cli.StringFlag{
					Name:        "host",
					Usage:       "The `HOST` to connect to",
					Destination: &Configuration.ClientOptions.Host,
					Value:       "localhost",
				},
			},
			Action: func(ctx *cli.Context) error {
				cfg := &Configuration
				c := client.NewClient()

				return c.Connect(cfg.ClientOptions.Host, cfg.Port)
			},
		},
	}

	return app
}
