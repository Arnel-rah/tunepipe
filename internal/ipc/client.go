package ipc

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"runtime"
	"sync"
)

type Client struct {
	conn    io.ReadWriteCloser
	scanner *bufio.Scanner
	mu      sync.Mutex
	reqID   uint64
}

type Command struct {
	Command   []interface{} `json:"command"`
	RequestID uint64        `json:"request_id,omitempty"`
}

type response struct {
	RequestID uint64          `json:"request_id"`
	Error     string          `json:"error"`
	Data      json.RawMessage `json:"data"`
}

func Connect(pipeName string) (*Client, error) {
	var conn io.ReadWriteCloser
	var err error

	if runtime.GOOS == "windows" {
		pipePath := `\\.\pipe\` + pipeName
		conn, err = os.OpenFile(pipePath, os.O_RDWR, 0)
	} else {
		socketPath := "/tmp/" + pipeName
		conn, err = net.Dial("unix", socketPath)
	}

	if err != nil {
		return nil, fmt.Errorf("impossible de se connecter a mpv (%s): %w", pipeName, err)
	}

	return &Client{
		conn:    conn,
		scanner: bufio.NewScanner(conn),
	}, nil
}

func (c *Client) SendExec(args ...interface{}) error {
	_, err := c.request(args...)
	return err
}

func (c *Client) GetProperty(name string) (json.RawMessage, error) {
	return c.request("get_property", name)
}

func (c *Client) request(args ...interface{}) (json.RawMessage, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.reqID++
	cmd := Command{
		Command:   args,
		RequestID: c.reqID,
	}

	data, err := json.Marshal(cmd)
	if err != nil {
		return nil, err
	}

	data = append(data, '\n')
	if _, err = c.conn.Write(data); err != nil {
		return nil, err
	}

	for c.scanner.Scan() {
		var result response
		if err := json.Unmarshal(c.scanner.Bytes(), &result); err != nil {
			continue
		}
		if result.RequestID != c.reqID {
			continue
		}
		if result.Error != "success" && result.Error != "" {
			return nil, fmt.Errorf("mpv command failed: %s", result.Error)
		}
		return result.Data, nil
	}
	if err := c.scanner.Err(); err != nil {
		return nil, err
	}
	return nil, io.EOF
}

func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
