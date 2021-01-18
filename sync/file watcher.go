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
	Root         string
	Recursive    bool
	FileCreated  func(file File) error
	FilesDeleted func(files []string) error
	Files        *[]File
	watcher      *fsnotify.Watcher
}

func (fw *FileWatcher) Watch(exit <-chan struct{}) error {
	var err error
	fw.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	defer fw.watcher.Close()

	err = fw.setupWatches()
	if err != nil {
		return err
	}

	go fw.listenForFsEvents()
	go fw.checkForDeletedFiles()

	<-exit

	return nil
}

func findDeletedFiles(files []File) []string {
	deleted := []string{}

	for _, f := range files {
		if _, err := os.Stat(f.Path); os.IsNotExist(err) {
			deleted = append(deleted, f.Path)
		}
	}

	return deleted
}

func (fw *FileWatcher) checkForDeletedFiles() {
	ticker := time.NewTicker(5 * time.Second)
	for {
		<-ticker.C

		deleted := findDeletedFiles(*fw.Files)
		if len(deleted) > 0 {
			fw.FilesDeleted(deleted)
		}
	}
}

func (fw *FileWatcher) listenForFsEvents() {
	for {
		select {
		case event := <-fw.watcher.Events:
			// watch for events
			fw.handleEvents(event)

		case err := <-fw.watcher.Errors:
			// watch for errors
			fmt.Printf("Watcher error: %#v", err)
		}
	}
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

func (fw *FileWatcher) handleEvents(e fsnotify.Event) {
	fileInfo, err := os.Stat(e.Name)

	// if the file does not exist
	if os.IsNotExist(err) {
		return
	}

	if err != nil {
		panic(err)
	}

	// We only care about files
	if fileInfo.IsDir() {
		return
	}

	if e.Op&fsnotify.Create == fsnotify.Create {
		fw.handleFileCreated(e)
		return
	}
}

func (fw *FileWatcher) handleFileCreated(e fsnotify.Event) {
	time.Sleep(1 * time.Second)

	file, err := GetFileInfo(fw.Root, e.Name)
	if err != nil {
		log.Printf("Error retreiving file: %v", err)
		return
	}

	fw.FileCreated(*file)
}

func (fw FileWatcher) IndexFiles(fileFound func(file File)) error {
	// files := []*File{}

	err := filepath.Walk(fw.Root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		file, err := GetFileInfo(fw.Root, path)
		if err != nil {
			return err
		}
		// files = append(files, file)
		fileFound(*file)

		return nil
	})

	if err != nil {
		return err
	}

	// return files, nil
	return nil
}
