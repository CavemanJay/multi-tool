package archive

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/CavemanJay/gogurt/config"
)

type Archiver struct {
	options config.ArchiveOptions
}

func NewArchiver(options config.ArchiveOptions) *Archiver {
	home, _ := os.UserHomeDir()

	for i, folder := range options.InFolders {
		options.InFolders[i] = strings.Replace(folder, "~", home, 1)
	}

	return &Archiver{
		options: options,
	}
}

func (a Archiver) Archive() error {
	pattern := regexp.MustCompile("Demos$")

	for _, folder := range a.options.InFolders {
		_, err := os.Stat(folder)
		if err != nil {
			return err
		}

		exePath, err := exec.LookPath("7z")
		if err != nil {
			return err
		}

		var archivePath string
		if pattern.MatchString(folder) {
			archiveName := strings.Replace(folder, "Demos", "RL Replays", 1)
			archivePath = filepath.Join(a.options.OutFolder, filepath.Base(archiveName))
		} else {
			archivePath = filepath.Join(a.options.OutFolder, filepath.Base(folder))
		}

		cmd := exec.Command(exePath, "u", archivePath, folder)
		cmd.Stdout = os.Stdout

		fmt.Println(cmd.String())
		err = cmd.Run()
		if err != nil {
			return err
		}

	}

	return nil
}
