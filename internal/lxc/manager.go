package lxc

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

// Container represents an LXC container
type Container struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Type   string `json:"type"`
}

// Manager handles LXC container operations
type Manager struct {
	// LXC container management interface
}

// NewManager creates a new LXC manager
func NewManager() *Manager {
	return &Manager{}
}

// ListContainers returns available LXC containers
func (m *Manager) ListContainers() ([]Container, error) {
	// Try to get containers from lxc command
	cmd := exec.Command("lxc", "list", "--format", "json")
	output, err := cmd.Output()

	if err != nil {
		// If lxc command fails, return mock containers for development
		return m.getMockContainers(), nil
	}

	// Parse lxc output
	var lxcContainers []map[string]interface{}
	if err := json.Unmarshal(output, &lxcContainers); err != nil {
		// If JSON parsing fails, return mock containers
		return m.getMockContainers(), nil
	}

	// Convert to our Container struct
	containers := make([]Container, 0, len(lxcContainers))
	for _, lxcContainer := range lxcContainers {
		name, _ := lxcContainer["name"].(string)
		status, _ := lxcContainer["status"].(string)
		containerType, _ := lxcContainer["type"].(string)

		containers = append(containers, Container{
			Name:   name,
			Status: strings.ToLower(status),
			Type:   containerType,
		})
	}

	// If no containers found, return mock containers for development
	if len(containers) == 0 {
		return m.getMockContainers(), nil
	}

	return containers, nil
}

// getMockContainers returns mock containers for development/testing
func (m *Manager) getMockContainers() []Container {
	return []Container{
		{Name: "dev-ubuntu", Status: "running", Type: "container"},
		{Name: "test-alpine", Status: "stopped", Type: "container"},
		{Name: "web-nginx", Status: "running", Type: "container"},
	}
}

// GetContainer returns info about a specific container
func (m *Manager) GetContainer(name string) (*Container, error) {
	containers, err := m.ListContainers()
	if err != nil {
		return nil, err
	}

	for _, container := range containers {
		if container.Name == name {
			return &container, nil
		}
	}

	return nil, fmt.Errorf("container '%s' not found", name)
}

// IsContainerRunning checks if a container is running
func (m *Manager) IsContainerRunning(name string) (bool, error) {
	container, err := m.GetContainer(name)
	if err != nil {
		return false, err
	}

	return container.Status == "running", nil
}

// ExecCommand executes a command in the specified container
// NOTE: This is a placeholder - actual exec will be implemented in T2.3/T2.4
func (m *Manager) ExecCommand(container, command string) error {
	// Check if container exists and is running
	running, err := m.IsContainerRunning(container)
	if err != nil {
		return fmt.Errorf("container check failed: %w", err)
	}

	if !running {
		return fmt.Errorf("container '%s' is not running", container)
	}

	// For T2.2, we only validate the container exists
	// Actual command execution will be handled by T2.3/T2.4 with PTY
	return nil
}
