package ytdlp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

const CacheDir = "cache"

var playcountsMu sync.Mutex

func EnsureCacheDir() error {
	if _, err := os.Stat(CacheDir); os.IsNotExist(err) {
		return os.MkdirAll(CacheDir, 0o755)
	}
	return nil
}

func CachedFilePath(id string) string {
	_ = EnsureCacheDir()
	exts := []string{"m4a", "mp3", "webm", "opus"}
	for _, e := range exts {
		p := filepath.Join(CacheDir, fmt.Sprintf("%s.%s", id, e))
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return ""
}

func DownloadToCache(id string) (string, error) {
	if err := EnsureCacheDir(); err != nil {
		return "", err
	}
	binary := getYTDLPBinary()
	outPattern := filepath.Join(CacheDir, "%(id)s.%(ext)s")
	args := []string{"-x", "--audio-format", "m4a", "--no-playlist", "-o", outPattern, fmt.Sprintf("https://www.youtube.com/watch?v=%s", id)}
	args = append(args, getCookiesOption()...)
	cmd := exec.Command(binary, args...)
	cmd.Stdout = nil
	cmd.Stderr = nil
	if err := cmd.Run(); err != nil {
		// best-effort: return error
		return "", fmt.Errorf("download error: %w", err)
	}
	if p := CachedFilePath(id); p != "" {
		return p, nil
	}
	return "", fmt.Errorf("downloaded but file not found")
}

func playcountsPath() string {
	_ = EnsureCacheDir()
	return filepath.Join(CacheDir, "playcounts.json")
}

func loadPlaycounts() (map[string]int, error) {
	m := make(map[string]int)
	p := playcountsPath()
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return m, nil
	}
	b, err := ioutil.ReadFile(p)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return m, nil
}

func savePlaycounts(m map[string]int) error {
	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	p := playcountsPath()
	tmp := p + ".tmp"
	if err := ioutil.WriteFile(tmp, b, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, p)
}

// IncrementPlayCount increments persistent play count for a track and returns the new count.
func IncrementPlayCount(id string) int {
	playcountsMu.Lock()
	defer playcountsMu.Unlock()
	m, _ := loadPlaycounts()
	m[id] = m[id] + 1
	_ = savePlaycounts(m)
	return m[id]
}

// EnsureDownloaded downloads the track if not already present. Returns path or error.
func EnsureDownloaded(id string) (string, error) {
	if p := CachedFilePath(id); p != "" {
		return p, nil
	}
	// avoid concurrent heavy downloads by a simple small sleep (best-effort)
	// real implementation could add per-id locks
	return DownloadToCache(id)
}
