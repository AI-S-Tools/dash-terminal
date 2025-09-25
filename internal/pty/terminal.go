package pty

import (
	"io"
	"os"
	"os/exec"
	"sync"

	"github.com/creack/pty"
)

// Terminal handles PTY (pseudo-terminal) operations
type Terminal struct {
	sessions map[string]*Session
	mutex    sync.RWMutex
}

// Session represents a PTY session
type Session struct {
	ID      string
	PTY     *os.File
	Command *exec.Cmd
	Running bool
	mutex   sync.RWMutex
}

// NewTerminal creates a new terminal manager
func NewTerminal() *Terminal {
	return &Terminal{
		sessions: make(map[string]*Session),
	}
}

// StartSession starts a new PTY session with given command
func (t *Terminal) StartSession(sessionID, containerName, command string) (*Session, error) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	// Create command - use lxc exec for container execution
	var cmd *exec.Cmd
	if containerName != "" {
		// Execute command inside LXC container
		cmd = exec.Command("lxc", "exec", containerName, "--", "bash", "-c", command)
	} else {
		// Execute command on host
		cmd = exec.Command("bash", "-c", command)
	}

	// Start PTY
	pttyFile, err := pty.Start(cmd)
	if err != nil {
		return nil, err
	}

	session := &Session{
		ID:      sessionID,
		PTY:     pttyFile,
		Command: cmd,
		Running: true,
	}

	t.sessions[sessionID] = session
	return session, nil
}

// GetSession returns a session by ID
func (t *Terminal) GetSession(sessionID string) (*Session, bool) {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	session, exists := t.sessions[sessionID]
	return session, exists
}

// Write writes data to the PTY session
func (s *Session) Write(data []byte) (int, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if !s.Running || s.PTY == nil {
		return 0, io.EOF
	}

	return s.PTY.Write(data)
}

// Read reads data from the PTY session
func (s *Session) Read(buffer []byte) (int, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if !s.Running || s.PTY == nil {
		return 0, io.EOF
	}

	return s.PTY.Read(buffer)
}

// Close closes the PTY session
func (s *Session) Close() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if !s.Running {
		return nil
	}

	s.Running = false

	if s.PTY != nil {
		s.PTY.Close()
	}

	if s.Command != nil && s.Command.Process != nil {
		s.Command.Process.Kill()
	}

	return nil
}

// Resize resizes the PTY terminal
func (s *Session) Resize(cols, rows int) error {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if !s.Running || s.PTY == nil {
		return io.EOF
	}

	return pty.Setsize(s.PTY, &pty.Winsize{
		Rows: uint16(rows),
		Cols: uint16(cols),
	})
}

// IsRunning returns whether the session is still running
func (s *Session) IsRunning() bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.Running
}
