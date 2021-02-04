package cli

import (
	"fmt"
	"os"
	"path/filepath"

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

	logFile, err := os.OpenFile(filepath.Join(cfg.AppDataFolder, "gogurt.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
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
	syncPath := filepath.Join(cfg.SyncFolder, chosenPlaylist.Name)

	toDownload := []*music.Video{}

	for _, video := range videos {
		fileName := fmt.Sprintf("%s.mp3", video.Title)
		_, err := os.Stat(filepath.Join(syncPath, fileName))
		if err != nil && os.IsNotExist(err) {
			toDownload = append(toDownload, &video)
		}
		err = nil
	}

	// err := music.DownloadVideo(videos[0], filepath.Join(cfg.SyncFolder, playlists[0].Name))
	// if err != nil {
	// 	log.Panic(err)
	// }

	return nil
}
