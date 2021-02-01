package archive

import (
	"fmt"
	"os"
	"os/exec"
	"path"
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
			archivePath = path.Join(a.options.OutFolder, path.Base(archiveName))
		} else {
			archivePath = path.Join(a.options.OutFolder, path.Base(folder))
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
