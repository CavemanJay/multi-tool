package sync

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
)

type FileWatcher struct {
	Root        string
	Recursive   bool
	FileCreated func(file *File) error
	watcher     *fsnotify.Watcher
}

func (fw *FileWatcher) Watch(cancel <-chan struct{}) error {
	var err error
	fw.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	// defer watcher.Close()

	err = fw.setupWatches()

	// done := make(chan bool, 1)

	go func() {
		for {
			select {
			case event := <-fw.watcher.Events:
				// watch for events
				fw.handleFileCreated(event)

			case err := <-fw.watcher.Errors:
				// watch for errors
				fmt.Println("ERROR", err)
			}
		}
	}()

	<-cancel

	return nil
}

func (fw *FileWatcher) setupWatches() error {
	// Walk the directory tree if we are in recursive mode
	if fw.Recursive {
		if err := filepath.Walk(fw.Root, fw.watchDir); err != nil {
			return err
		}
	} else {
		if err := fw.watcher.Add(fw.Root); err != nil {
			return err
		}
	}

	return nil
}

func (fw *FileWatcher) watchDir(path string, fi os.FileInfo, err error) error {
	// since fsnotify can watch all the files in a directory, watchers only need
	// to be added to each nested directory
	if fi.Mode().IsDir() {
		return fw.watcher.Add(path)
	}

	return nil
}

func (fw *FileWatcher) handleFileCreated(e fsnotify.Event) {
	fileInfo, err := os.Stat(e.Name)

	// if the file does not exist
	if os.IsNotExist(err) {
		return
	}

	if err != nil {
		panic(err)
	}

	if !fileInfo.IsDir() && e.Op&fsnotify.Create == fsnotify.Create {
		time.Sleep(1 * time.Second)

		file, err := GetFileInfo(e.Name)
		if err != nil {
			log.Printf("Error retreiving file: %v", err)
			return
		}

		fw.FileCreated(file)
	}
}
