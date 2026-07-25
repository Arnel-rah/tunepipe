package ytdlp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Track struct {
	ID       string  `json:"id"`
	Title    string  `json:"title"`
	Uploader string  `json:"uploader"`
	Duration float64 `json:"duration"`
}

func getYTDLPBinary() string {
	if _, err := os.Stat(".\\yt-dlp.exe"); err == nil {
		return ".\\yt-dlp.exe"
	}
	return "yt-dlp"
}

func Search(query string, limit int) ([]Track, error) {
	binary := getYTDLPBinary()
	searchQuery := fmt.Sprintf("ytsearch%d:%s", limit, query)
	cmd := exec.Command(binary,
		"--flat-playlist",
		"-j",
		"--no-warnings",
		"--force-ipv4",
		"--no-cache-dir",
		searchQuery,
	)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		errMsg := strings.TrimSpace(stderr.String())
		if errMsg != "" {
			var errResponse map[string]interface{}
			if json.Unmarshal([]byte(errMsg), &errResponse) == nil {
				if msg, ok := errResponse["error"].(string); ok {
					return nil, fmt.Errorf("yt-dlp error: %s", msg)
				}
			}
			return nil, fmt.Errorf("%s", errMsg)
		}
		return nil, fmt.Errorf("erreur recherche yt-dlp: %w", err)
	}

	var tracks []Track
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		var t Track
		if err := json.Unmarshal([]byte(line), &t); err == nil {
			tracks = append(tracks, t)
		}
	}

	return tracks, nil
}

func FetchDirectURL(videoID string) (string, error) {
	binary := getYTDLPBinary()
	url := fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoID)
	cmd := exec.Command(binary,
		"-f", "ba/ba*",
		"-g",
		"--no-warnings",
		"--force-ipv4",
		"--no-cache-dir",
		url,
	)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		errMsg := strings.TrimSpace(stderr.String())
		if errMsg != "" {
			return "", fmt.Errorf("%s", errMsg)
		}
		return "", fmt.Errorf("impossible de recuperer le flux: %w", err)
	}

	return strings.TrimSpace(out.String()), nil
}
