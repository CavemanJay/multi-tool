package cli

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/CavemanJay/multi-tool/archive"
	"github.com/CavemanJay/multi-tool/client"
	"github.com/CavemanJay/multi-tool/music"
	"github.com/CavemanJay/multi-tool/server"
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

	logFile, err := os.OpenFile(filepath.Join(cfg.AppDataFolder, "multi-tool.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer logFile.Close()
	initLogger(logFile)

	if cfg.Append != "" {
		cfg.SyncFolder = filepath.Join(cfg.SyncFolder, cfg.Append)
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
	cfg := &configuration
	if err := handleConfig(); err != nil {
		return err
	}

	// initLogger(nil)
	writeConfig(cfg)

	client := music.NewYoutubeClient(cfg.MusicOptions.SecretsFile)
	playlists := client.Playlists()
	var chosenPlaylist *music.PlayList

	chosen := cfg.MusicOptions.PlaylistName

	// Playlist not specified. Ask what playlist the user wants to download
	if len(chosen) <= 0 {
		chosen = getPlaylist(&playlists)
	}

	if chosen == "quit" {
		os.Exit(0)
	}

	for _, pl := range playlists {
		if pl.Name == chosen {
			chosenPlaylist = &pl
			break
		}
	}

	// If the playlist wasn't found
	if chosenPlaylist == nil {
		return fmt.Errorf("Playlist '%s' does not exist", chosen)
	}

	videos := client.Videos(chosenPlaylist)
	syncPath := filepath.Join(cfg.SyncFolder, "Music", chosenPlaylist.Name)

	toDownload := []*music.Video{}

	var limit int
	if cfg.MusicOptions.Limit == 0 {
		limit = len(videos)
	} else if cfg.MusicOptions.Limit >= len(videos) {
		limit = len(videos)
	} else {
		limit = cfg.MusicOptions.Limit
	}

	// TODO: Optimize
	count := 0
	for i := 0; count < limit && i < len(videos); i++ {
		fileName := music.GetFileName(&videos[i])
		toCheck := filepath.Join(syncPath, fileName)
		_, err := os.Stat(toCheck)
		if err != nil && os.IsNotExist(err) {
			toDownload = append(toDownload, &videos[i])
			count++
		}
		err = nil
	}

	for _, v := range toDownload {
		delay := rand.Intn(4-1) + 1
		err := music.DownloadVideo(v, syncPath)
		if err != nil {
			return err
		}

		time.Sleep(time.Duration(delay) * time.Second)
	}

	return nil
}
