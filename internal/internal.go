package internal

import (
	"NewPhotoWeb/log"
	"bytes"
	"fmt"
	"os/exec"
)

func CreateThumbnailFromVideo(video []byte) []byte {
	var buf bytes.Buffer
	const (
		width  = 500
		height = 500
	)

	cmd := exec.Command("ffmpeg", "-i", "pipe:0", "-vframes", "1", "-s", fmt.Sprintf("%dx%d", width, height), "-f", "singlejpeg", "-")
	cmd.Stdout = &buf
	cmd.Stdin = bytes.NewReader(video)

	if err := cmd.Start(); err != nil {
		log.Logger.Fatalln(err)
	}

	if err := cmd.Wait(); err != nil {
		log.Logger.Fatalln(err)
	}

	return buf.Bytes()
}
