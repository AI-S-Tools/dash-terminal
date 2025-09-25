#!/bin/bash

# Create remaining numbered tasks for dash-terminal

# P2 tasks
echo "Creating P2 tasks..."
dppm task create t2-2-lxc-manager --project dash-terminal --phase p2-core-functionality --title "T2.2: LXC Manager" --description "Go interface for LXC container operations and lifecycle" --priority high

dppm task create t2-3-terminal-wrapper --project dash-terminal --phase p2-core-functionality --title "T2.3: Terminal Wrapper" --description "JavaScript wrapper for xterm.js with WebSocket integration" --priority medium

dppm task create t2-4-native-tmux-manager --project dash-terminal --phase p2-core-functionality --title "T2.4: Native Tmux Manager" --description "Go tmux session/window/pane management with layout parsing" --priority high

# P3 tasks
echo "Creating P3 tasks..."
dppm task create t3-1-session-tabs --project dash-terminal --phase p3-mobile-polish --title "T3.1: Session Tab Component" --description "Interactive session tabs with swipe gestures and mobile optimization" --priority medium

dppm task create t3-2-window-tabs --project dash-terminal --phase p3-mobile-polish --title "T3.2: Window Tab Component" --description "Window tab switching with creation/deletion and tmux sync" --priority medium

dppm task create t3-3-pane-grid-system --project dash-terminal --phase p3-mobile-polish --title "T3.3: Pane Grid System" --description "Dynamic pane layout rendering with tmux layout parsing and zoom" --priority high

dppm task create t3-4-touch-gestures --project dash-terminal --phase p3-mobile-polish --title "T3.4: Touch Gesture Handler" --description "Mobile touch gestures for pane operations and navigation" --priority high

echo "All numbered tasks created!"
dppm status project dash-terminal