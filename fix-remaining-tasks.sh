#!/bin/bash

# Fix remaining 8 tasks (T2.1 - T3.4) with full anti-scope creep

DPPM_BASE="/home/ubuntu/Dropbox/project-management/projects/dash-terminal"

echo "🔧 Fixing remaining T2.1 - T3.4 tasks..."

# T2.1: WebSocket Handler
cat > "$DPPM_BASE/phases/p2-core-functionality/tasks/t2-1-websocket-handler.yaml" << 'EOF'
id: t2-1-websocket-handler
title: 'T2.1: WebSocket Handler'
project_id: dash-terminal
phase_id: p2-core-functionality
status: todo
priority: high
reporter: dppm-user
created: "2025-09-25"
updated: "2025-09-25"
dependency_ids: ["t1-1-project-setup", "t1-4-tmux-websocket-protocol"]
description: |
  ## T2.1: WebSocket Handler - SERVER CONNECTION ONLY

  **⚠️ STOP! CHECK DEPENDENCIES FIRST:**
  ```bash
  # T1.1 AND T1.4 MUST BE COMPLETE:
  dppm task show t1-1-project-setup --project dash-terminal | grep "status: done" || { echo "ERROR: T1.1 not done"; exit 1; }
  dppm task show t1-4-tmux-websocket-protocol --project dash-terminal | grep "status: done" || { echo "ERROR: T1.4 not done"; exit 1; }
  ls -la internal/websocket/protocol.go || { echo "ERROR: T1.4 incomplete - no message types"; exit 1; }
  ```

  **📋 SUMMARY FROM PREVIOUS TASKS:**
  - ✅ **T1.1**: Go project with HTTP server + WebSocket handler stub
  - ✅ **T1.2**: PWA foundation (not needed for this task)
  - ✅ **T1.3**: CSS tabs (not needed for this task)
  - ✅ **T1.4**: Message types defined in protocol.go
  - 📍 **You are here**: T2.1 - Replace WebSocket stub with real server

  **📂 FILE STATUS CHECK:**
  ```bash
  grep "not implemented" internal/websocket/handler.go || { echo "WARNING: Already implemented"; exit 1; }
  ```

  **✏️ FILE OPERATION MODE: REFACTOR EXISTING**
  - REFACTOR: Replace stub in internal/websocket/handler.go
  - Use 'Read' first to see current stub
  - Use 'Edit' to replace "not implemented" with real WebSocket server
  - DO NOT create new files

  **🔒 PROTECTED FILES:**
  - internal/websocket/protocol.go (T1.4's work - do not modify)
  - cmd/server/main.go (only add WebSocket endpoint)
  - internal/lxc/, internal/tmux/ (other tasks)
  - web/ files (frontend)

  **🚫 FORBIDDEN - DO NOT ADD:**
  - ❌ NO actual tmux command execution (T2.4's job)
  - ❌ NO LXC container operations (T2.2's job)
  - ❌ NO frontend JavaScript (T2.3's job)
  - ❌ NO authentication/authorization
  - ❌ NO session persistence/database
  - ❌ NO HTTP endpoints beyond WebSocket upgrade
  - ❌ NO business logic - just message routing

  **✅ VERIFICATION:**
  ```bash
  go build ./cmd/server || echo "FAIL: Build broken"
  # Manual test: WebSocket connection should accept but return "not implemented" for commands
  ```

  **🎯 SUCCESS CRITERIA:**
  - WebSocket server accepts connections ✓
  - Routes messages by type ✓
  - Returns "not implemented" for actual commands ✓
  - Does NOT execute tmux/LXC commands ✓

EOF

# T2.2: LXC Manager
cat > "$DPPM_BASE/phases/p2-core-functionality/tasks/t2-2-lxc-manager.yaml" << 'EOF'
id: t2-2-lxc-manager
title: 'T2.2: LXC Manager'
project_id: dash-terminal
phase_id: p2-core-functionality
status: todo
priority: high
reporter: dppm-user
created: "2025-09-25"
updated: "2025-09-25"
dependency_ids: ["t1-1-project-setup"]
description: |
  ## T2.2: LXC Manager - INTERFACE IMPLEMENTATION ONLY

  **⚠️ STOP! CHECK DEPENDENCIES FIRST:**
  ```bash
  dppm task show t1-1-project-setup --project dash-terminal | grep "status: done" || { echo "ERROR: T1.1 not done"; exit 1; }
  ls -la internal/lxc/manager.go || { echo "ERROR: T1.1 incomplete - no LXC stub"; exit 1; }
  ```

  **📋 SUMMARY FROM PREVIOUS TASKS:**
  - ✅ **T1.1**: Created internal/lxc/manager.go stub with interface
  - 📍 **You are here**: T2.2 - Implement LXC container discovery only

  **📂 FILE STATUS CHECK:**
  ```bash
  grep "not implemented" internal/lxc/manager.go || { echo "WARNING: Already implemented"; exit 1; }
  ```

  **✏️ FILE OPERATION MODE: REFACTOR EXISTING**
  - REFACTOR: Replace stub in internal/lxc/manager.go
  - Implement ONLY ListContainers() and basic container info
  - NO container lifecycle operations yet

  **🚫 FORBIDDEN - DO NOT ADD:**
  - ❌ NO container creation/deletion (out of scope)
  - ❌ NO container start/stop operations
  - ❌ NO file operations in containers
  - ❌ NO network configuration
  - ❌ NO WebSocket integration (T2.1's job)
  - ❌ NO tmux session management (T2.4's job)
  - ❌ ONLY read-only operations: list, info, status

  **🎯 SUCCESS CRITERIA:**
  - Lists available LXC containers ✓
  - Returns container status (running/stopped) ✓
  - NO modification operations ✓

EOF

echo "T2.1 and T2.2 fixed. Creating remaining T2.3, T2.4, T3.1-T3.4..."

# Continue with remaining tasks...
echo "Script needs continuation for T2.3-T3.4. Run phase 2..."