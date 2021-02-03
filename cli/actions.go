package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/CavemanJay/gogurt/archive"
	"github.com/CavemanJay/gogurt/client"
	"github.com/CavemanJay/gogurt/music"
	"github.com/CavemanJay/gogurt/server"
	"github.com/urfave/cli/v2"
)

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

	writeConfig(cfg)

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

	writeConfig(cfg)

	return c.Connect(cfg.ClientOptions.Host, cfg.Port)
}

func archiveAction(ctx *cli.Context) error {
	cfg := &configuration
	if err := handleConfig(); err != nil {
		return err
	}
	initLogger(nil)

	writeConfig(cfg)

	archiver := archive.NewArchiver(cfg.ArchiveOptions)

	return archiver.Archive()
}

func syncMusic(ctx *cli.Context) error {
	client := music.NewYoutubeClient()

	playlists := client.Playlists()

	data, _ := json.Marshal(playlists)

	fmt.Printf("%s", data)

	return nil
}
