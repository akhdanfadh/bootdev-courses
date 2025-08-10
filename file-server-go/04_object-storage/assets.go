package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"os/exec"
)

type FFProbeStream struct {
	CodecType string `json:"codec_type"`
	Width     int    `json:"width"`  // only in video, 0 if not present
	Height    int    `json:"height"` // only in video, 0 if not present
}

func (cfg apiConfig) ensureAssetsDir() error {
	if _, err := os.Stat(cfg.assetsRoot); os.IsNotExist(err) {
		return os.Mkdir(cfg.assetsRoot, 0755)
	}
	return nil
}

func getVideoAspectRatio(filePath string) (string, error) {
	// Execute ffprobe command to get video stream information
	cmd := exec.Command("ffprobe", "-v", "error", "-print_format", "json", "-show_streams", filePath)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	// Get height and width of the first video stream
	var response struct {
		Streams []FFProbeStream `json:"streams"`
	}
	if err := json.Unmarshal(output, &response); err != nil {
		return "", err
	}
	width, height := 0, 0
	for _, stream := range response.Streams {
		if stream.CodecType == "video" {
			width = stream.Width
			height = stream.Height
			break
		}
	}
	if width == 0 || height == 0 {
		return "", fmt.Errorf("no valid video stream")
	}

	// Handle aspect ratio calculation
	const landscapeDefault = 16.0 / 9.0
	const landscapeTolerance = 0.1
	const portraitDefault = 9.0 / 16.0
	const portrainTolerance = 0.05
	aspectRatio := float64(width) / float64(height)
	if aspectRatio >= 1 && math.Abs(aspectRatio-landscapeDefault) <= landscapeTolerance {
		return "16:9", nil
	} else if aspectRatio < 1 && math.Abs(aspectRatio-portraitDefault) <= portrainTolerance {
		return "9:16", nil
	} else {
		return "other", nil
	}
}
