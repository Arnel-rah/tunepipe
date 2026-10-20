package player

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"time"

	"github.com/Arnel-rah/tunepipe/internal/ipc"
)

func (e *Engine) TimePosition() (float64, error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.ipcClient == nil {
		return 0, fmt.Errorf("ipc client not connected")
	}
	data, err := e.ipcClient.GetProperty("time-pos")
	if err != nil {
		return 0, err
	}
	var position *float64
	if err := json.Unmarshal(data, &position); err != nil {
		return 0, err
	}
	if position == nil {
		return 0, fmt.Errorf("mpv has no active time position")
	}
	return *position, nil
}

type Engine struct {
	cmd       *exec.Cmd
	ipcClient *ipc.Client
	pipeName  string
	mu        sync.Mutex
}

func getMPVBinary() string {
	if _, err := os.Stat(".\\mpv.exe"); err == nil {
		return ".\\mpv.exe"
	}
	return "mpv"
}

func NewEngine(pipeName string, ytDlpPath string) (*Engine, error) {
	binary := getMPVBinary()

	var pipePath string
	if runtime.GOOS == "windows" {
		pipePath = fmt.Sprintf("\\\\.\\pipe\\%s", pipeName)
	} else {
		pipePath = "/tmp/" + pipeName
	}

	args := []string{
		"--idle",
		"--no-video",
		"--no-terminal",
		"--volume=100",
		"--mute=no",
		"--audio-device=auto",
		"--audio-channels=stereo",
		fmt.Sprintf("--input-ipc-server=%s", pipePath),
		"--ytdl-format=bestaudio[abr<=64]/bestaudio[abr<=96]/bestaudio[abr<=128]/bestaudio",
		"--cache=yes",
		"--demuxer-max-bytes=8MiB",
		"--demuxer-max-back-bytes=1MiB",
		"--demuxer-readahead-secs=20",
	}
	if runtime.GOOS == "linux" {
		args = append(args, "--ao=pipewire")
	} else {
		args = append(args, "--ao=auto")
	}

	if ytDlpPath != "" {
		args = append(args, fmt.Sprintf("--script-opts=ytdl_hook-ytdl_path=%s", ytDlpPath))
	}

	cmd := exec.Command(binary, args...)

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("impossible de demarrer mpv: %w", err)
	}

	var client *ipc.Client
	var err error

	for i := 0; i < 15; i++ {
		time.Sleep(200 * time.Millisecond)
		client, err = ipc.Connect(pipeName)
		if err == nil {
			break
		}
	}

	if err != nil {
		_ = cmd.Process.Kill()
		return nil, fmt.Errorf("echec connexion IPC mpv: %w", err)
	}

	return &Engine{
		cmd:       cmd,
		ipcClient: client,
		pipeName:  pipeName,
	}, nil
}

func (e *Engine) PlayURL(url string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.ipcClient == nil {
		return fmt.Errorf("ipc client not connected")
	}

	if err := e.ipcClient.SendExec("stop"); err != nil {
		return err
	}

	if err := e.ipcClient.SendExec("loadfile", url, "replace"); err != nil {
		return err
	}

	time.Sleep(200 * time.Millisecond)

	if err := e.ipcClient.SendExec("set_property", "volume", 100); err != nil {
		return err
	}
	if err := e.ipcClient.SendExec("set_property", "mute", false); err != nil {
		return err
	}
	for i := 0; i < 3; i++ {
		if err := e.ipcClient.SendExec("set_property", "pause", false); err == nil {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}

	return nil
}

func (e *Engine) TogglePause() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.ipcClient == nil {
		return fmt.Errorf("ipc client not connected")
	}

	return e.ipcClient.SendExec("cycle", "pause")
}

func (e *Engine) Pause() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.ipcClient == nil {
		return fmt.Errorf("ipc client not connected")
	}

	return e.ipcClient.SendExec("set", "pause", "yes")
}

func (e *Engine) Resume() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.ipcClient == nil {
		return fmt.Errorf("ipc client not connected")
	}

	for i := 0; i < 3; i++ {
		if err := e.ipcClient.SendExec("set", "pause", "no"); err == nil {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}

	return nil
}

func (e *Engine) SetVolume(volume int) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.ipcClient == nil {
		return fmt.Errorf("ipc client not connected")
	}
	if volume < 0 {
		volume = 0
	}
	if volume > 100 {
		volume = 100
	}

	if err := e.ipcClient.SendExec("set_property", "volume", volume); err != nil {
		return err
	}

	return e.ipcClient.SendExec("set_property", "mute", volume == 0)
}

func (e *Engine) Stop() {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e == nil {
		return
	}
	if e.ipcClient != nil {
		_ = e.ipcClient.SendExec("quit")
		_ = e.ipcClient.Close()
		e.ipcClient = nil
	}
	if e.cmd != nil && e.cmd.Process != nil {
		_ = e.cmd.Process.Kill()
		_ = e.cmd.Wait()
	}
}
