package cli

import (
	"io"
	"log"
	"os"
	"path"

	"github.com/CavemanJay/gogurt/config"
	"github.com/op/go-logging"
)

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

	stdOutBackend := logging.NewLogBackend(os.Stdout, "", 0)
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

func writeConfig(cfg *config.Config) {
	config.WriteConfig(path.Join(cfg.AppDataFolder, "last_run.json"), cfg)
}
