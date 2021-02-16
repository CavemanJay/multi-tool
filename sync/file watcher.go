package sync

import (
	"os"
	"path/filepath"
	"time"

	"github.com/karrick/godirwalk"
	"github.com/op/go-logging"

	"github.com/fsnotify/fsnotify"
)

var log = logging.MustGetLogger("multi-tool")

type FileWatcher struct {
	Root        string
	Recursive   bool
	FileCreated func(file File) error
	FileDeleted func(path string) error
	Files       *[]File
	watcher     *fsnotify.Watcher
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

	<-exit

	return nil
}

func (fw *FileWatcher) listenForFsEvents() {
	for {
		select {
		case event := <-fw.watcher.Events:
			// watch for events
			fw.handleEvents(event)

		case err := <-fw.watcher.Errors:
			// watch for errors
			log.Errorf("Watcher error: %#v", err)
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

	// if the file does not exist (file is deleted)
	if os.IsNotExist(err) {
		log.Debugf("File deleted: %s", e.Name)
		fw.FileDeleted(getFilePathRelativeToRoot(fw.Root, e.Name))
		return
	}

	if err != nil {
		panic(err)
	}

	if fileInfo.IsDir() {
		err := fw.watcher.Add(e.Name)
		if err != nil {
			log.Error(err)
		}
		return
	}

	if e.Op&fsnotify.Create == fsnotify.Create {
		log.Debugf("File created: %s", e.Name)
		fw.handleFileCreated(e)
		return
	}
}

func (fw *FileWatcher) handleFileCreated(e fsnotify.Event) {
	time.Sleep(1 * time.Second)

	file, err := GetFileInfo(fw.Root, e.Name)
	if err != nil {
		log.Errorf("Error retreiving file: %v", err)
		return
	}

	fw.FileCreated(*file)
}

func (fw FileWatcher) IndexFiles(fileFound func(file File)) error {
	count := 0
	err := godirwalk.Walk(fw.Root, &godirwalk.Options{
		Callback: func(path string, de *godirwalk.Dirent) error {
			if !de.IsDir() {
				count++
			}
			return nil
		},
	})
	if err != nil {
		log.Error(err)
	}
	// bar := pb.StartNew(count)
	// defer bar.Finish()

	// indexed := 0
	return godirwalk.Walk(fw.Root, &godirwalk.Options{
		Callback: func(path string, de *godirwalk.Dirent) error {
			if de.IsDir() {
				return nil
			}

			file, err := GetFileInfo(fw.Root, path)
			if err != nil {
				return err
			}

			fileFound(*file)
			// bar.Increment()
			// fmt.Printf("\rIndexed file %d of %d", indexed, count)
			// indexed++

			return nil
		},
		Unsorted: true,
	})
}
