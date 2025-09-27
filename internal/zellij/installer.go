package zellij

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// EnsureZellijInstalled checks if Zellij is installed and installs it if not
func (d *Daemon) EnsureZellijInstalled() error {
	// Check if Zellij is already available
	if d.isZellijAvailable() {
		return nil
	}

	// Try to install Zellij
	return d.installZellij()
}

// isZellijAvailable checks if Zellij command is available
func (d *Daemon) isZellijAvailable() bool {
	var cmd *exec.Cmd
	if d.container == "" || d.container == "host" {
		cmd = exec.Command("zellij", "--version")
	} else {
		cmd = exec.Command("lxc", "exec", d.container, "--", "zellij", "--version")
	}

	err := cmd.Run()
	return err == nil
}

// installZellij installs Zellij using the official installer
func (d *Daemon) installZellij() error {
	if d.container == "" || d.container == "host" {
		return d.installZellijOnHost()
	} else {
		return d.installZellijInContainer()
	}
}

// installZellijOnHost installs Zellij on the host system
func (d *Daemon) installZellijOnHost() error {
	// Create installation directory
	installDir := filepath.Join(os.Getenv("HOME"), ".local", "bin")
	if err := os.MkdirAll(installDir, 0755); err != nil {
		return fmt.Errorf("failed to create install directory: %v", err)
	}

	// Download and install Zellij
	var downloadURL string
	switch runtime.GOARCH {
	case "amd64":
		downloadURL = "https://github.com/zellij-org/zellij/releases/latest/download/zellij-x86_64-unknown-linux-musl.tar.gz"
	case "arm64":
		downloadURL = "https://github.com/zellij-org/zellij/releases/latest/download/zellij-aarch64-unknown-linux-musl.tar.gz"
	default:
		return fmt.Errorf("unsupported architecture: %s", runtime.GOARCH)
	}

	// Download and extract
	installScript := fmt.Sprintf(`
		set -e
		cd /tmp
		wget -O zellij.tar.gz "%s"
		tar -xzf zellij.tar.gz
		mv zellij "%s/"
		chmod +x "%s/zellij"
		rm -f zellij.tar.gz
	`, downloadURL, installDir, installDir)

	cmd := exec.Command("bash", "-c", installScript)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install Zellij: %v", err)
	}

	return nil
}

// installZellijInContainer installs Zellij inside an LXC container
func (d *Daemon) installZellijInContainer() error {
	// First ensure the container has required tools
	setupCmd := exec.Command("lxc", "exec", d.container, "--", "bash", "-c", `
		apt-get update -qq || apk update || yum update -y || true
		apt-get install -y wget tar || apk add wget tar || yum install -y wget tar || true
		mkdir -p /usr/local/bin
	`)

	if err := setupCmd.Run(); err != nil {
		return fmt.Errorf("failed to setup container dependencies: %v", err)
	}

	// Determine architecture
	archCmd := exec.Command("lxc", "exec", d.container, "--", "uname", "-m")
	archOutput, err := archCmd.Output()
	if err != nil {
		return fmt.Errorf("failed to determine container architecture: %v", err)
	}

	var downloadURL string
	arch := string(archOutput)
	if arch == "x86_64\n" || arch == "x86_64" {
		downloadURL = "https://github.com/zellij-org/zellij/releases/latest/download/zellij-x86_64-unknown-linux-musl.tar.gz"
	} else if arch == "aarch64\n" || arch == "aarch64" {
		downloadURL = "https://github.com/zellij-org/zellij/releases/latest/download/zellij-aarch64-unknown-linux-musl.tar.gz"
	} else {
		return fmt.Errorf("unsupported container architecture: %s", arch)
	}

	// Install Zellij in container
	installScript := fmt.Sprintf(`
		set -e
		cd /tmp
		wget -O zellij.tar.gz "%s"
		tar -xzf zellij.tar.gz
		mv zellij /usr/local/bin/
		chmod +x /usr/local/bin/zellij
		rm -f zellij.tar.gz
	`, downloadURL)

	installCmd := exec.Command("lxc", "exec", d.container, "--", "bash", "-c", installScript)
	if err := installCmd.Run(); err != nil {
		return fmt.Errorf("failed to install Zellij in container: %v", err)
	}

	return nil
}

// GetZellijVersion returns the installed Zellij version
func (d *Daemon) GetZellijVersion() (string, error) {
	var cmd *exec.Cmd
	if d.container == "" || d.container == "host" {
		cmd = exec.Command("zellij", "--version")
	} else {
		cmd = exec.Command("lxc", "exec", d.container, "--", "zellij", "--version")
	}

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}