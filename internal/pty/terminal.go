package pty

import (
	"io"
	"os"
	"os/exec"
	"sync"

	"github.com/creack/pty"
)

// Session represents a PTY session
type Session struct {
	pty     *os.File
	command *exec.Cmd
	running bool
	mutex   sync.RWMutex
}

// NewSession starts a new PTY session with a shell inside the specified container.
func NewSession(containerName string) (*Session, error) {
	var cmd *exec.Cmd
	if containerName == "" || containerName == "host" {
		cmd = exec.Command("bash")
	} else {
		cmd = exec.Command("lxc", "exec", containerName, "--", "bash")
	}

	pttyFile, err := pty.Start(cmd)
	if err != nil {
		return nil, err
	}

	session := &Session{
		pty:     pttyFile,
		command: cmd,
		running: true,
	}

	return session, nil
}

// Write writes data to the PTY session
func (s *Session) Write(data []byte) (int, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if !s.running || s.pty == nil {
		return 0, io.EOF
	}

	return s.pty.Write(data)
}

// Read reads data from the PTY session
func (s *Session) Read(buffer []byte) (int, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if !s.running || s.pty == nil {
		return 0, io.EOF
	}

	return s.pty.Read(buffer)
}

// Close closes the PTY session
func (s *Session) Close() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if !s.running {
		return nil
	}

	s.running = false

	if s.pty != nil {
		s.pty.Close()
	}

	if s.command != nil && s.command.Process != nil {
		s.command.Process.Kill()
	}

	return nil
}

// Resize resizes the PTY terminal
func (s *Session) Resize(width, height int) error {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if !s.running || s.pty == nil {
		return io.EOF
	}

	return pty.Setsize(s.pty, &pty.Winsize{
		Rows: uint16(height),
		Cols: uint16(width),
	})
}