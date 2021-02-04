package music

import (
	"log"
	"os"
	"os/exec"
	"path"
)

const template = "%(title)s.%(format)s"

var exePath string

func init() {
	var err error
	exePath, err = exec.LookPath("youtube-dl")
	if err != nil {
		log.Panic("youtube-dl executable not found")
	}
}

// func quote(str string) string {
// 	return fmt.Sprintf("'%s'", str)
// }

func DownloadVideo(video Video, outputRoot string) error {
	cmd := exec.Command(
		exePath,
		"--ignore-config",
		"-w",
		"--extract-audio",
		"--audio-format",
		"mp3",
		"-o",
		path.Join(outputRoot, template), video.Link())

	// fmt.Println(cmd.String())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
