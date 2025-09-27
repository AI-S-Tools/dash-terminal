package zellij

import (
	"os"
	"os/exec"

	"github.com/creack/pty"
)

// startWithPTY starts a command with a PTY and returns the PTY file
func startWithPTY(cmd *exec.Cmd) (*os.File, error) {
	// Start the command with a PTY
	ptyFile, err := pty.Start(cmd)
	if err != nil {
		return nil, err
	}

	return ptyFile, nil
}

// ResizePTY resizes the PTY to the specified dimensions
func ResizePTY(ptyFile *os.File, width, height int) error {
	return pty.Setsize(ptyFile, &pty.Winsize{
		Rows: uint16(height),
		Cols: uint16(width),
	})
}