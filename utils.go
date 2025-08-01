package main

import (
	"bytes"
	"encoding/json"
	"os/exec"
)

type ffprobeStream struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type ffprobe struct {
	Streams []ffprobeStream `json:"streams"`
}

func getVideoAspectRatio(filePath string) (string, error) {
	ffprobeCmd := exec.Command(
		"ffprobe",
		"-v",
		"error",
		"-print_format",
		"json",
		"-show_streams",
		filePath,
	)
	ffprobeBuffer := bytes.Buffer{}
	ffprobeCmd.Stdout = &ffprobeBuffer
	err := ffprobeCmd.Run()
	if err != nil {
		return "", err
	}

	var ffprobeResult ffprobe
	err = json.Unmarshal(ffprobeBuffer.Bytes(), &ffprobeResult)
	if err != nil {
		return "", err
	}

	landscapeRatio := float32(16) / 9
	portraitRatio := float32(9) / 16

	tolerance := float32(1) / (1 << 5)

	videoStream := ffprobeResult.Streams[0]
	videoRatio := float32(videoStream.Width) / float32(videoStream.Height)

	if -tolerance <= videoRatio-landscapeRatio && videoRatio-landscapeRatio <= tolerance {
		return "landscape", nil
	} else if -tolerance <= videoRatio-portraitRatio && videoRatio-portraitRatio <= tolerance {
		return "portrait", nil
	}
	return "other", nil
}

func processVideoForFastStart(filePath string) (string, error) {
	outputFilePath := filePath + ".processing"

	ffprobeCmd := exec.Command(
		"ffmpeg",
		"-i",
		filePath,
		"-c",
		"copy",
		"-movflags",
		"faststart",
		"-f",
		"mp4",
		outputFilePath,
	)
	err := ffprobeCmd.Run()
	if err != nil {
		return "", err
	}

	return outputFilePath, nil
}
