package tmux

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// Manager handles tmux session operations
type Manager struct {
	containerName string
}

// Session represents a tmux session
type Session struct {
	Name    string   `json:"name"`
	Windows []Window `json:"windows"`
}

// Window represents a tmux window
type Window struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Panes []Pane `json:"panes"`
}

// Pane represents a tmux pane
type Pane struct {
	ID     int  `json:"id"`
	Active bool `json:"active"`
	Width  int  `json:"width"`
	Height int  `json:"height"`
}

// NewManager creates a new tmux manager
func NewManager(containerName string) *Manager {
	return &Manager{
		containerName: containerName,
	}
}

// ListSessions returns all tmux sessions
func (m *Manager) ListSessions() ([]Session, error) {
	// Get list of tmux sessions
	cmd := m.buildCommand("list-sessions", "-F", "#{session_name}")
	output, err := cmd.Output()

	if err != nil {
		// No sessions exist - return empty list
		return []Session{}, nil
	}

	sessionNames := strings.Split(strings.TrimSpace(string(output)), "\n")
	sessions := make([]Session, 0, len(sessionNames))

	for _, name := range sessionNames {
		if name == "" {
			continue
		}

		session := Session{Name: name}

		// Get windows for this session
		windows, err := m.listWindows(name)
		if err == nil {
			session.Windows = windows
		}

		sessions = append(sessions, session)
	}

	return sessions, nil
}

// CreateSession creates a new tmux session
func (m *Manager) CreateSession(name string) error {
	cmd := m.buildCommand("new-session", "-d", "-s", name)
	return cmd.Run()
}

// listWindows lists all windows in a session
func (m *Manager) listWindows(sessionName string) ([]Window, error) {
	cmd := m.buildCommand("list-windows", "-t", sessionName, "-F", "#{window_index}:#{window_name}")
	output, err := cmd.Output()

	if err != nil {
		return []Window{}, nil
	}

	windowLines := strings.Split(strings.TrimSpace(string(output)), "\n")
	windows := make([]Window, 0, len(windowLines))

	for _, line := range windowLines {
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		id, err := strconv.Atoi(parts[0])
		if err != nil {
			continue
		}

		window := Window{
			ID:   id,
			Name: parts[1],
		}

		// Get panes for this window
		panes, err := m.listPanes(sessionName, id)
		if err == nil {
			window.Panes = panes
		}

		windows = append(windows, window)
	}

	return windows, nil
}

// listPanes lists all panes in a window
func (m *Manager) listPanes(sessionName string, windowID int) ([]Pane, error) {
	target := fmt.Sprintf("%s:%d", sessionName, windowID)
	cmd := m.buildCommand("list-panes", "-t", target, "-F", "#{pane_index}:#{?pane_active,1,0}:#{pane_width}:#{pane_height}")
	output, err := cmd.Output()

	if err != nil {
		return []Pane{}, nil
	}

	paneLines := strings.Split(strings.TrimSpace(string(output)), "\n")
	panes := make([]Pane, 0, len(paneLines))

	for _, line := range paneLines {
		if line == "" {
			continue
		}

		parts := strings.Split(line, ":")
		if len(parts) != 4 {
			continue
		}

		id, err := strconv.Atoi(parts[0])
		if err != nil {
			continue
		}

		active := parts[1] == "1"
		width, _ := strconv.Atoi(parts[2])
		height, _ := strconv.Atoi(parts[3])

		panes = append(panes, Pane{
			ID:     id,
			Active: active,
			Width:  width,
			Height: height,
		})
	}

	return panes, nil
}

// CreateWindow creates a new window in a session
func (m *Manager) CreateWindow(sessionName, windowName string) error {
	cmd := m.buildCommand("new-window", "-t", sessionName, "-n", windowName)
	return cmd.Run()
}

// CreatePane creates a new pane in a window
func (m *Manager) CreatePane(sessionName string, windowID int) error {
	target := fmt.Sprintf("%s:%d", sessionName, windowID)
	cmd := m.buildCommand("split-window", "-t", target)
	return cmd.Run()
}

// SelectSession switches to a session
func (m *Manager) SelectSession(sessionName string) error {
	cmd := m.buildCommand("switch-client", "-t", sessionName)
	return cmd.Run()
}

// SelectWindow switches to a window
func (m *Manager) SelectWindow(sessionName string, windowID int) error {
	target := fmt.Sprintf("%s:%d", sessionName, windowID)
	cmd := m.buildCommand("select-window", "-t", target)
	return cmd.Run()
}

// SelectPane switches to a pane
func (m *Manager) SelectPane(sessionName string, windowID, paneID int) error {
	target := fmt.Sprintf("%s:%d.%d", sessionName, windowID, paneID)
	cmd := m.buildCommand("select-pane", "-t", target)
	return cmd.Run()
}

// SendKeys sends keystrokes to a pane
func (m *Manager) SendKeys(sessionName string, windowID, paneID int, keys string) error {
	target := fmt.Sprintf("%s:%d.%d", sessionName, windowID, paneID)
	cmd := m.buildCommand("send-keys", "-t", target, keys)
	return cmd.Run()
}

// buildCommand builds a tmux command, optionally executing in a container
func (m *Manager) buildCommand(args ...string) *exec.Cmd {
	if m.containerName != "" {
		// Execute tmux inside LXC container
		fullArgs := append([]string{"exec", m.containerName, "--", "tmux"}, args...)
		return exec.Command("lxc", fullArgs...)
	} else {
		// Execute tmux on host
		return exec.Command("tmux", args...)
	}
}
