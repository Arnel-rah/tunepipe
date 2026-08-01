package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/Arnel-rah/tunepipe/internal/player"
	"github.com/Arnel-rah/tunepipe/internal/ui"
	"github.com/Arnel-rah/tunepipe/internal/ytdlp"

	tea "github.com/charmbracelet/bubbletea"
)

func ytDlpURL() string {
	switch runtime.GOOS {
	case "windows":
		return "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp.exe"
	case "darwin":
		return "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp_macos"
	default:
		return "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp"
	}
}

func localBinDir() (string, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(base, "tunepipe", "bin")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	return dir, nil
}

func downloadYtDlp() (string, error) {
	dir, err := localBinDir()
	if err != nil {
		return "", err
	}

	name := "yt-dlp"
	if runtime.GOOS == "windows" {
		name = "yt-dlp.exe"
	}
	dest := filepath.Join(dir, name)

	fmt.Println("yt-dlp introuvable, téléchargement en cours...")
	resp, err := http.Get(ytDlpURL())
	if err != nil {
		return "", fmt.Errorf("téléchargement échoué: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("téléchargement échoué: statut HTTP %d", resp.StatusCode)
	}

	out, err := os.Create(dest)
	if err != nil {
		return "", err
	}
	defer out.Close()

	if _, err := io.Copy(out, resp.Body); err != nil {
		return "", err
	}

	if runtime.GOOS != "windows" {
		if err := os.Chmod(dest, 0o755); err != nil {
			return "", err
		}
	}

	fmt.Println("yt-dlp installé avec succès.")
	return dest, nil
}

func ensureYtDlp() (string, error) {
	if path, err := exec.LookPath("yt-dlp"); err == nil {
		if exec.Command(path, "--version").Run() == nil {
			return path, nil
		}
	}

	dir, err := localBinDir()
	if err == nil {
		name := "yt-dlp"
		if runtime.GOOS == "windows" {
			name = "yt-dlp.exe"
		}
		localPath := filepath.Join(dir, name)
		if _, statErr := os.Stat(localPath); statErr == nil {
			if exec.Command(localPath, "--version").Run() == nil {
				return localPath, nil
			}
		}
	}

	return downloadYtDlp()
}

func mpvInstallHint() string {
	switch runtime.GOOS {
	case "windows":
		return "winget install mpv-player.mpv"
	case "darwin":
		return "brew install mpv"
	default:
		return "sudo apt install mpv   (ou l'équivalent pour ta distro)"
	}
}

func ensureMpv() error {
	path, err := exec.LookPath("mpv")
	if err == nil {
		if exec.Command(path, "--version").Run() == nil {
			return nil
		}
	}

	return fmt.Errorf(
		"mpv est introuvable ou ne fonctionne pas correctement.\n  Installe-le avec : %s\n  Puis relance tunepipe.",
		mpvInstallHint(),
	)
}

func main() {
	ytDlpPath, err := ensureYtDlp()
	if err != nil {
		fmt.Printf("Impossible de préparer yt-dlp : %v\n", err)
		fmt.Println("Installe-le manuellement depuis : https://github.com/yt-dlp/yt-dlp/releases")
		os.Exit(1)
	}
	ytdlp.SetBinary(ytDlpPath)
	if err := ensureMpv(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	pipeName := "tunepipe-mpv"
	eng, err := player.NewEngine(pipeName, ytDlpPath)
	if err != nil {
		fmt.Printf("Erreur au demarrage de mpv: %v\n", err)
		os.Exit(1)
	}
	defer eng.Stop()

	p := tea.NewProgram(ui.NewModel(eng), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Erreur TUI: %v\n", err)
		os.Exit(1)
	}
}
