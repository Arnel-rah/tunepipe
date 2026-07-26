package player

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"

	"tunepipe/internal/ipc"
)

type Engine struct {
	cmd       *exec.Cmd
	ipcClient *ipc.Client
	pipeName  string
}

func getMPVBinary() string {
	if _, err := os.Stat(".\\mpv.exe"); err == nil {
		return ".\\mpv.exe"
	}
	return "mpv"
}

func NewEngine(pipeName string) (*Engine, error) {
	binary := getMPVBinary()

	var pipePath string
	if runtime.GOOS == "windows" {
		pipePath = fmt.Sprintf("\\\\.\\pipe\\%s", pipeName)
	} else {
		pipePath = "/tmp/" + pipeName
	}

	cmd := exec.Command(binary,
		"--idle",
		"--no-video",
		fmt.Sprintf("--input-ipc-server=%s", pipePath),
	)

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
	if err := e.ipcClient.SendExec("loadfile", url, "replace"); err != nil {
		return err
	}
	
	return e.ipcClient.SendExec("set_property", "pause", false)
}

func (e *Engine) TogglePause() error {
	return e.ipcClient.SendExec("cycle", "pause")
}

func (e *Engine) Stop() {
	if e == nil {
		return
	}
	if e.ipcClient != nil {
		_ = e.ipcClient.SendExec("quit")
		_ = e.ipcClient.Close()
	}
	if e.cmd != nil && e.cmd.Process != nil {
		_ = e.cmd.Process.Kill()
		_ = e.cmd.Wait()
	}
}
