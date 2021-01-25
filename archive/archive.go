package archive

import (
	"github.com/CavemanJay/gogurt/config"
	"github.com/kjk/lzmadec"
)

type Archiver struct {
	config.ArchiveOptions
}

func NewArchiver(options config.ArchiveOptions) *Archiver {
	return &Archiver{
		ArchiveOptions: options,
	}
}

func getPaths(entries *[]lzmadec.Entry) []string {
	var paths = []string{}

	for _, entry := range *entries {
		paths = append(paths, entry.Path)
	}

	return paths
}
