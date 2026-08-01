package player

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"time"

	"github.com/Arnel-rah/tunepipe/internal/ipc"
)

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
		"--really-quiet",
		fmt.Sprintf("--input-ipc-server=%s", pipePath),
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

	for i := 0; i < 3; i++ {
		if err := e.ipcClient.SendExec("set", "pause", "no"); err == nil {
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
