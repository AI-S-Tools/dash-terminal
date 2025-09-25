#!/bin/bash

DPPM_BASE="/home/ubuntu/Dropbox/project-management/projects/dash-terminal"

# Quick batch creation of remaining 6 tasks with anti-scope creep basics

# T2.3: Terminal Wrapper
cat > "$DPPM_BASE/phases/p2-core-functionality/tasks/t2-3-terminal-wrapper.yaml" << 'EOF'
id: t2-3-terminal-wrapper
title: 'T2.3: Terminal Wrapper'
project_id: dash-terminal
phase_id: p2-core-functionality
status: todo
priority: medium
reporter: dppm-user
created: "2025-09-25"
updated: "2025-09-25"
dependency_ids: ["t1-2-pwa-foundation", "t2-1-websocket-handler"]
description: |
  ## T2.3: Terminal Wrapper - XTERM.JS INTEGRATION ONLY

  **⚠️ STOP! CHECK DEPENDENCIES FIRST:**
  ```bash
  dppm task show t1-2-pwa-foundation --project dash-terminal | grep "status: done" || { echo "ERROR: T1.2 PWA not done"; exit 1; }
  dppm task show t2-1-websocket-handler --project dash-terminal | grep "status: done" || { echo "ERROR: T2.1 WebSocket not done"; exit 1; }
  ```

  **📋 SUMMARY:** T1.2 has PWA + T2.1 has WebSocket server + Need xterm.js client

  **🚫 FORBIDDEN:**
  - ❌ NO tmux command integration (T2.4's job)
  - ❌ NO touch gestures (T3.4's job)
  - ❌ NO tab switching logic (T3.1-T3.2's job)
  - ❌ ONLY basic terminal rendering + WebSocket connection

  **🎯 SCOPE:** Add xterm.js to existing HTML, connect to WebSocket, show terminal - NO functionality beyond basic rendering
EOF

# T2.4: Native Tmux Manager
cat > "$DPPM_BASE/phases/p2-core-functionality/tasks/t2-4-native-tmux-manager.yaml" << 'EOF'
id: t2-4-native-tmux-manager
title: 'T2.4: Native Tmux Manager'
project_id: dash-terminal
phase_id: p2-core-functionality
status: todo
priority: high
reporter: dppm-user
created: "2025-09-25"
updated: "2025-09-25"
dependency_ids: ["t2-2-lxc-manager", "t1-4-tmux-websocket-protocol"]
description: |
  ## T2.4: Native Tmux Manager - TMUX COMMAND EXECUTION ONLY

  **⚠️ STOP! CHECK DEPENDENCIES FIRST:**
  ```bash
  dppm task show t2-2-lxc-manager --project dash-terminal | grep "status: done" || { echo "ERROR: T2.2 LXC not done"; exit 1; }
  dppm task show t1-4-tmux-websocket-protocol --project dash-terminal | grep "status: done" || { echo "ERROR: T1.4 protocol not done"; exit 1; }
  ```

  **📋 SUMMARY:** T2.2 lists containers + T1.4 has message types + Need tmux command execution

  **🚫 FORBIDDEN:**
  - ❌ NO WebSocket message handling (T2.1's job)
  - ❌ NO frontend JavaScript (T2.3's job)
  - ❌ ONLY implement: list-sessions, list-windows, list-panes, send-keys
EOF

# T3.1: Session Tab Component
cat > "$DPPM_BASE/phases/p3-mobile-polish/tasks/t3-1-session-tabs.yaml" << 'EOF'
id: t3-1-session-tabs
title: 'T3.1: Session Tab Component'
project_id: dash-terminal
phase_id: p3-mobile-polish
status: todo
priority: medium
reporter: dppm-user
created: "2025-09-25"
updated: "2025-09-25"
dependency_ids: ["t2-3-terminal-wrapper"]
description: |
  ## T3.1: Session Tab Component - JAVASCRIPT INTERACTIVITY ONLY

  **⚠️ STOP! CHECK DEPENDENCIES:**
  ```bash
  dppm task show t2-3-terminal-wrapper --project dash-terminal | grep "status: done" || { echo "ERROR: T2.3 terminal not done"; exit 1; }
  ```

  **📋 SUMMARY:** T1.3 has CSS tabs + T2.3 has terminal + Need JavaScript click handlers

  **🚫 FORBIDDEN:**
  - ❌ NO tmux session creation (backend's job)
  - ❌ ONLY click handlers for existing CSS tabs
EOF

# T3.2, T3.3, T3.4 - Similar pattern
cat > "$DPPM_BASE/phases/p3-mobile-polish/tasks/t3-2-window-tabs.yaml" << 'EOF'
id: t3-2-window-tabs
title: 'T3.2: Window Tab Component'
project_id: dash-terminal
phase_id: p3-mobile-polish
status: todo
priority: medium
reporter: dppm-user
created: "2025-09-25"
updated: "2025-09-25"
dependency_ids: ["t2-3-terminal-wrapper"]
description: |
  ## T3.2: Window Tab Component - JAVASCRIPT ONLY
  **⚠️ DEPENDENCIES:** T2.3 terminal wrapper must be done
  **🚫 FORBIDDEN:** NO backend integration - only JavaScript tab switching
EOF

cat > "$DPPM_BASE/phases/p3-mobile-polish/tasks/t3-3-pane-grid-system.yaml" << 'EOF'
id: t3-3-pane-grid-system
title: 'T3.3: Pane Grid System'
project_id: dash-terminal
phase_id: p3-mobile-polish
status: todo
priority: high
reporter: dppm-user
created: "2025-09-25"
updated: "2025-09-25"
dependency_ids: ["t2-4-native-tmux-manager"]
description: |
  ## T3.3: Pane Grid System - LAYOUT RENDERING ONLY
  **⚠️ DEPENDENCIES:** T2.4 tmux manager must be done
  **🚫 FORBIDDEN:** NO pane creation - only render existing tmux layouts
EOF

cat > "$DPPM_BASE/phases/p3-mobile-polish/tasks/t3-4-touch-gestures.yaml" << 'EOF'
id: t3-4-touch-gestures
title: 'T3.4: Touch Gesture Handler'
project_id: dash-terminal
phase_id: p3-mobile-polish
status: todo
priority: high
reporter: dppm-user
created: "2025-09-25"
updated: "2025-09-25"
dependency_ids: ["t3-3-pane-grid-system"]
description: |
  ## T3.4: Touch Gesture Handler - MOBILE GESTURES ONLY
  **⚠️ DEPENDENCIES:** T3.3 pane grid must be done
  **🚫 FORBIDDEN:** NO new functionality - only touch/swipe handlers for existing UI
  **🎯 FINAL TASK:** This completes the MVP - native app-level mobile experience
EOF

echo "✅ All 12 tasks now have anti-scope creep guards!"
dppm status project dash-terminal