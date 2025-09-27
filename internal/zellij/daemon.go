package zellij

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"sync"
	"time"
)

// Daemon manages Zellij sessions for a container
type Daemon struct {
	container string
	sessions  map[string]*Session
	mu        sync.RWMutex
}

// Session represents a Zellij session
type Session struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Created   time.Time `json:"created"`
	Active    bool      `json:"active"`
	Container string    `json:"container"`
	Process   *exec.Cmd `json:"-"`
	mu        sync.RWMutex
}

// NewDaemon creates a new Zellij daemon
func NewDaemon(container string) *Daemon {
	daemon := &Daemon{
		container: container,
		sessions:  make(map[string]*Session),
	}

	// Ensure Zellij is installed
	if err := daemon.EnsureZellijInstalled(); err != nil {
		fmt.Printf("Warning: Failed to ensure Zellij installation: %v\n", err)
		fmt.Println("Zellij features may not work properly")
	} else {
		version, err := daemon.GetZellijVersion()
		if err == nil {
			fmt.Printf("Zellij available: %s", version)
		}
	}

	return daemon
}

// CreateSession creates a new Zellij session
func (d *Daemon) CreateSession(name string) (*Session, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	// Generate session ID
	sessionID := fmt.Sprintf("%s-%d", name, time.Now().Unix())

	// Create Zellij session using wrapper script for better PTY compatibility
	var cmd *exec.Cmd
	if d.container == "" || d.container == "host" {
		// Run Zellij wrapper on host
		cmd = exec.Command("./zellij-wrapper.sh", sessionID)
	} else {
		// Run Zellij wrapper in LXC container
		cmd = exec.Command("lxc", "exec", d.container, "--", "/usr/local/bin/zellij-wrapper.sh", sessionID)
	}

	session := &Session{
		ID:        sessionID,
		Name:      name,
		Created:   time.Now(),
		Active:    true,
		Container: d.container,
		Process:   cmd,
	}

	d.sessions[sessionID] = session
	return session, nil
}

// GetSession retrieves a session by ID
func (d *Daemon) GetSession(sessionID string) (*Session, bool) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	session, exists := d.sessions[sessionID]
	return session, exists
}

// ListSessions returns all sessions
func (d *Daemon) ListSessions() []*Session {
	d.mu.RLock()
	defer d.mu.RUnlock()

	sessions := make([]*Session, 0, len(d.sessions))
	for _, session := range d.sessions {
		sessions = append(sessions, session)
	}
	return sessions
}

// GetContainer returns the container name
func (d *Daemon) GetContainer() string {
	return d.container
}

// DeleteSession removes and terminates a session
func (d *Daemon) DeleteSession(sessionID string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	session, exists := d.sessions[sessionID]
	if !exists {
		return fmt.Errorf("session %s not found", sessionID)
	}

	// Kill Zellij session
	if err := d.killZellijSession(sessionID); err != nil {
		return fmt.Errorf("failed to kill session: %v", err)
	}

	// Mark session as inactive
	session.mu.Lock()
	session.Active = false
	session.mu.Unlock()

	delete(d.sessions, sessionID)
	return nil
}

// killZellijSession kills a Zellij session
func (d *Daemon) killZellijSession(sessionID string) error {
	var cmd *exec.Cmd
	if d.container == "" || d.container == "host" {
		cmd = exec.Command("zellij", "kill-session", sessionID)
	} else {
		cmd = exec.Command("lxc", "exec", d.container, "--", "zellij", "kill-session", sessionID)
	}

	return cmd.Run()
}

// AttachToSession attaches to an existing session via PTY
func (d *Daemon) AttachToSession(sessionID string) (*os.File, error) {
	session, exists := d.GetSession(sessionID)
	if !exists {
		return nil, fmt.Errorf("session %s not found", sessionID)
	}

	// Start the session process with PTY if not already started
	if session.Process != nil {
		ptyFile, err := startWithPTY(session.Process)
		if err != nil {
			return nil, fmt.Errorf("failed to start session: %v", err)
		}
		return ptyFile, nil
	}

	return nil, fmt.Errorf("session process not available")
}

// GetSessionInfo returns detailed info about a session using Zellij commands
func (d *Daemon) GetSessionInfo(sessionID string) (map[string]interface{}, error) {
	var cmd *exec.Cmd
	if d.container == "" || d.container == "host" {
		cmd = exec.Command("zellij", "list-sessions")
	} else {
		cmd = exec.Command("lxc", "exec", d.container, "--", "zellij", "list-sessions")
	}

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get session info: %v", err)
	}

	// Parse Zellij output (basic implementation)
	info := make(map[string]interface{})
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(output, len(output))

	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			info["raw_output"] = string(output)
			break
		}
	}

	return info, nil
}