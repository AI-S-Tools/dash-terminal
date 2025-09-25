package lxc

// Manager handles LXC container operations
// This is a stub file - implementation in T2.2
type Manager struct {
	// LXC container management interface
}

// ListContainers returns available LXC containers
// Implementation in T2.2
func (m *Manager) ListContainers() ([]string, error) {
	return nil, nil
}

// ExecCommand executes a command in the specified container
// Implementation in T2.2
func (m *Manager) ExecCommand(container, command string) error {
	return nil
}
