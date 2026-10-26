package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"os/exec"
)

func getVideoAspectRatio(path string) (string, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-print_format", "json", "-show_streams", path)
	var output bytes.Buffer
	cmd.Stdout = &output

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	size := struct {
		Streams []struct {
			Width  int `json:"width"`
			Height int `json:"height"`
		} `json:"streams"`
	}{}

	err = json.Unmarshal(output.Bytes(), &size)
	if err != nil {
		return "", err
	}

	if len(size.Streams) < 1 {
		return "", errors.New("no streams in video")
	}

	width := size.Streams[0].Width
	height := size.Streams[0].Height

	if (width * 9) / (height * 16) == 1 {
		return "16:9", nil
	}

	if (width * 16) / (height * 9) == 1 {
		return "9:16", nil
	}

	return "other", nil
}
