#!/bin/bash

# Fix dppm task dependencies according to tasklist.md

DPPM_PATH="/home/ubuntu/Dropbox/project-management/projects/dash-terminal"

# Foundation Phase - T1 tasks
echo "Fixing Foundation Phase dependencies..."

# T1.4 depends on T1.1
sed -i '/^updated:/a dependency_ids: ["project-setup"]' "$DPPM_PATH/phases/foundation/tasks/tmux-websocket-protocol.yaml" 2>/dev/null

# Core Functionality Phase - T2 tasks
echo "Fixing Core Functionality Phase dependencies..."

# T2.1 depends on T1.1 and T1.4
sed -i '/^updated:/a dependency_ids: ["project-setup", "tmux-websocket-protocol"]' "$DPPM_PATH/phases/core-functionality/tasks/websocket-handler.yaml"

# T2.2 depends on T1.1
sed -i '/^updated:/a dependency_ids: ["project-setup"]' "$DPPM_PATH/phases/core-functionality/tasks/lxc-manager.yaml"

# T2.3 depends on T1.2 and T2.1
sed -i '/^updated:/a dependency_ids: ["pwa-foundation", "websocket-handler"]' "$DPPM_PATH/phases/core-functionality/tasks/terminal-wrapper.yaml"

# T2.4 depends on T2.2 and T1.4
sed -i '/^updated:/a dependency_ids: ["lxc-manager", "tmux-websocket-protocol"]' "$DPPM_PATH/phases/core-functionality/tasks/native-tmux-manager.yaml"

# Mobile Polish Phase - T3 tasks
echo "Fixing Mobile Polish Phase dependencies..."

# T3.1 depends on T2.3
sed -i '/^updated:/a dependency_ids: ["terminal-wrapper"]' "$DPPM_PATH/phases/mobile-polish/tasks/session-tabs.yaml"

# T3.2 depends on T2.3
sed -i '/^updated:/a dependency_ids: ["terminal-wrapper"]' "$DPPM_PATH/phases/mobile-polish/tasks/window-tabs.yaml"

# T3.3 depends on T2.4
sed -i '/^updated:/a dependency_ids: ["native-tmux-manager"]' "$DPPM_PATH/phases/mobile-polish/tasks/pane-grid-system.yaml"

# T3.4 depends on T3.3
sed -i '/^updated:/a dependency_ids: ["pane-grid-system"]' "$DPPM_PATH/phases/mobile-polish/tasks/touch-gestures.yaml"

echo "Dependencies fixed!"
dppm status dependencies --project dash-terminal