package ytdlp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type Track struct {
	ID       string  `json:"id"`
	Title    string  `json:"title"`
	Uploader string  `json:"uploader"`
	Duration float64 `json:"duration"`
}

var configuredBinary string

func SetBinary(path string) {
	configuredBinary = path
}

func getYTDLPBinary() string {
	if configuredBinary != "" {
		return configuredBinary
	}
	if _, err := os.Stat(".\\yt-dlp.exe"); err == nil {
		return ".\\yt-dlp.exe"
	}
	return "yt-dlp"
}

func getCookiesOption() []string {
	binary := getYTDLPBinary()
	if runtime.GOOS == "windows" {
		if _, err := os.Stat("cookies.txt"); err == nil {
			return []string{"--cookies", "cookies.txt"}
		}
		browsers := []string{"firefox", "chrome", "edge", "brave"}
		for _, b := range browsers {
			check := exec.Command(binary, "--cookies-from-browser", b, "--cookies", "test.txt")
			if err := check.Run(); err == nil {
				os.Remove("test.txt")
				return []string{"--cookies-from-browser", b}
			}
		}
		return []string{}
	}
	return []string{"--cookies-from-browser", "firefox"}
}

func Search(query string, limit int) ([]Track, error) {
	binary := getYTDLPBinary()
	searchQuery := fmt.Sprintf("ytsearch%d:%s", limit, query)

	args := []string{
		"--flat-playlist",
		"-j",
		"--no-warnings",
		"--force-ipv4",
		"--remote-components", "ejs:github",
		"--user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:128.0) Gecko/20100101 Firefox/128.0",
	}

	args = append(args, getCookiesOption()...)
	args = append(args, searchQuery)

	cmd := exec.Command(binary, args...)

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

func SearchWithCache(query string, limit int) ([]Track, error) {
	cache := GetSearchCache()
	cacheKey := fmt.Sprintf("%d:%s", limit, query)
	if results, ok := cache.Get(cacheKey); ok {
		return results, nil
	}

	results, err := Search(query, limit)
	if err != nil {
		return nil, err
	}

	cache.Set(cacheKey, results)
	return results, nil
}

// formatSelectors est ordonné pour privilégier les débits audio bas
// (économie de data), en évitant tout fallback vers un format vidéo+audio
// combiné qui ferait exploser la consommation réseau.
func formatSelectors() []string {
	return []string{
		"bestaudio[ext=m4a][abr<=64]",
		"bestaudio[ext=webm][abr<=64]",
		"bestaudio[acodec=opus][abr<=96]",
		"bestaudio[abr<=96]",
		"bestaudio[abr<=128]",
		"bestaudio",
		"worstaudio", // dernier recours: pire format AUDIO, jamais de vidéo
	}
}

func tryFormats(binary, url string, baseArgs []string) (string, error) {
	var lastErr error
	var lastOutput string

	for _, selector := range formatSelectors() {
		args := append([]string{"-f", selector}, baseArgs...)
		args = append(args, url)

		cmd := exec.Command(binary, args...)

		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr

		if err := cmd.Run(); err != nil {
			errMsg := strings.TrimSpace(stderr.String())
			if errMsg != "" && !strings.Contains(errMsg, "Requested format is not available") &&
				!strings.Contains(errMsg, "No video formats found") {
				lastErr = fmt.Errorf("%s", errMsg)
			}
			continue
		}

		directURL := strings.TrimSpace(out.String())
		if directURL != "" && strings.HasPrefix(directURL, "http") {
			return directURL, nil
		}
		if directURL != "" {
			lastOutput = directURL
		}
	}

	if lastOutput != "" {
		return lastOutput, nil
	}

	if lastErr != nil {
		return "", lastErr
	}
	return "", fmt.Errorf("Tsisy format dispo")
}

func FetchDirectURL(videoID string) (string, error) {
	binary := getYTDLPBinary()
	url := fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoID)

	baseArgs := []string{
		"-g",
		"--no-warnings",
		"--force-ipv4",
		"--remote-components", "ejs:github",
		"--user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:128.0) Gecko/20100101 Firefox/128.0",
	}

	directURL, err := tryFormats(binary, url, baseArgs)
	if err == nil && directURL != "" {
		return directURL, nil
	}

	withCookies := append(append([]string{}, baseArgs...), getCookiesOption()...)
	directURL, err = tryFormats(binary, url, withCookies)
	if err == nil && directURL != "" {
		return directURL, nil
	}

	return "", fmt.Errorf("mbola tsisy format iany: %w", err)
}