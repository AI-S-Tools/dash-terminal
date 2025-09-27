#!/bin/bash

# Zellij wrapper for PTY integration
# This script makes Zellij work better with PTY from Go

SESSION_NAME="$1"
shift  # Remove session name from arguments

# Set environment for better PTY compatibility
export TERM=xterm-256color
export COLORTERM=truecolor
export ZELLIJ_AUTO_ATTACH=false

# Start Zellij with the specified session
exec zellij --session "$SESSION_NAME" "$@"