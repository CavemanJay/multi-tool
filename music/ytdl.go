package music

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

const template = "%(title)s.%(format)s"

var exePath string

func init() {
	exePath, _ = exec.LookPath("youtube-dl")
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
		"-o",
		filepath.Join(outputRoot, template), video.Link())

	fmt.Printf("Downloading track: %s\n", video.Title)
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr

	return cmd.Run()
}
