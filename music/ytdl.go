package music

import (
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"strings"
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

func GetFileName(video *Video) string {
	fileName := strings.ReplaceAll(video.Title, "/", "_")
	fileName = strings.ReplaceAll(fileName, "**OUT ON SPOTIFY**", "_OUT ON SPOTIFY")
	fileName = strings.ReplaceAll(fileName, "|", "_")
	// re := regexp.MustCompile(`\**$`)
	// fileName = re.ReplaceAllString(fileName, "")
	fileName = strings.ReplaceAll(fileName, "*", "_")

	fileName = fmt.Sprintf("%s.mp3", fileName)
	return fileName
}

func DownloadVideo(video *Video, outputRoot string) error {
	cmd := exec.Command(
		exePath,
		"--ignore-config",
		"-w",
		"--extract-audio",
		"--audio-format",
		"mp3",
		"--audio-quality",
		"9",
		"-o",
		filepath.Join(outputRoot, template), video.Link())

	fmt.Printf("Downloading track: %s\n", video.Title)
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr

	return cmd.Run()
}
