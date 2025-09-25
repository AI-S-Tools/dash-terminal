#!/bin/bash

# Quick fix for T1.1 YAML with proper indentation

DPPM_BASE="/home/ubuntu/Dropbox/project-management/projects/dash-terminal"

cat > "$DPPM_BASE/phases/p1-foundation/tasks/t1-1-project-setup.yaml" << 'EOF'
id: t1-1-project-setup
title: 'T1.1: Project Boilerplate Setup'
project_id: dash-terminal
phase_id: p1-foundation
status: todo
priority: critical
reporter: dppm-user
created: "2025-09-25"
updated: "2025-09-25"
dependency_ids: []
description: |
  ## T1.1: Project Boilerplate Setup - GO PTY SERVER FOUNDATION ONLY

  **🎯 PROJEKT OVERSIGT (GENTAGES I HVER TASK FOR FULD KONTEKST):**

  Vi bygger **LXC Terminal PWA** - en mobil terminal applikation til administration af LXC containers.
  Dette er ikke bare en web-app, men en installérbar PWA der føles som en native iOS/Android app.

  **🏗️ ARKITEKTUR OVERVIEW:**
  - **Backend:** Go PTY server med WebSocket kommunikation på port :8080
  - **Frontend:** Vanilla JavaScript + xterm.js terminal emulator (IKKE React/Vue for performance)
  - **Integration:** Native tmux session management med tre-niveau hierarki: Sessions → Windows → Panes
  - **Platform:** PWA med service worker, installérbar på mobile enheder

  **📱 MOBILE-FIRST DESIGN PRINCIPPER:**
  - iOS/Android safe area support (notch handling)
  - Touch-optimized interface med gesture support
  - 60fps scrolling og smooth animations
  - Sub-50ms response time for terminal commands
  - Bundle size <2MB for hurtig loading på mobile netværk
  - Dark theme matching terminal æstetik

  **🔧 TECH STACK DETALJER:**
  ```
  Backend (Go):
  ├── gorilla/websocket for WebSocket connections
  ├── creack/pty for PTY (pseudo-terminal) handling
  ├── Native LXC integration via lxc commands
  └── tmux native integration (sessions/windows/panes)

  Frontend (Vanilla JS):
  ├── xterm.js terminal emulator library
  ├── Native Web APIs (ServiceWorker, WebSocket, PWA Manifest)
  ├── CSS Grid for responsive pane layouts
  └── Touch/gesture handling for mobile interaction

  Development:
  ├── Test-driven development (TDD) required
  ├── git-air for commits (custom commit tool)
  ├── DPPM for project management
  └── Anti-scope creep system for AI development
  ```

  **🎯 PERFORMANCE TARGETS (KRITISKE REQUIREMENTS):**
  - Terminal response: <50ms fra command til display
  - Scrolling: Konstant 60fps på mobile devices
  - Bundle size: <2MB total (critical for mobile)
  - Memory usage: <50MB på mobile browser
  - Battery impact: Minimal drain ved kontinuerlig brug

  **🚫 SCOPE BOUNDARIES (HVAD VI IKKE BYGGER):**
  - ❌ IKKE multi-user support (single user PWA)
  - ❌ IKKE file upload/download features
  - ❌ IKKE authentication system
  - ❌ IKKE multi-server support (single LXC host)
  - ❌ IKKE container statistics/monitoring dashboards
  - ❌ IKKE themes eller customization settings
  - ❌ IKKE logging/analytics systems
  - ❌ IKKE collaboration features

  **📋 PROJEKT PHASES OVERVIEW:**
  - **Phase 1 (Foundation):** T1.1-T1.4 - Basic project structure, PWA, CSS framework, protocols
  - **Phase 2 (Core):** T2.1-T2.4 - WebSocket server, LXC integration, terminal functionality
  - **Phase 3 (Polish):** T3.1-T3.4 - Mobile UI components, gestures, final optimizations

  **⚠️ STOP! OBLIGATORISK DEPENDENCY CHECK:**
  ```bash
  # NO DEPENDENCIES - This is the foundation task
  echo "T1.1 is the first task - no dependencies to check"
  echo "Ready to create complete Go project structure"
  ```

  **📋 HANDOFF SUMMARY FRA TIDLIGERE TASKS:**
  **PREVIOUS TASKS:** None - this is the foundation task
  **YOU ARE HERE:** T1.1 - Create complete Go project structure from scratch

  **📂 FILE STATUS CHECK - KØR DETTE FØRST:**
  ```bash
  # CHECK IF THIS TASK IS ALREADY DONE
  echo "Checking if T1.1 project setup is already complete..."
  ls -la cmd/server/main.go 2>/dev/null && echo "WARNING: main.go EXISTS - task may be done" || echo "OK: main.go does not exist - ready to start"
  ls -la go.mod 2>/dev/null && echo "WARNING: go.mod EXISTS - project may be initialized" || echo "OK: go.mod missing - ready to initialize"
  ```

  **✏️ FILE OPERATION MODE: CREATE NEW PROJECT**
  - This task CREATES NEW project structure from scratch
  - Use 'Write' tool to create all new files
  - Project follows Go standards with cmd/, internal/, web/ structure
  - LXC integration setup (not Docker)
  - Basic HTTP server on :8080 serving web/ directory

  **🔒 PROTECTED FILES - MUST NOT MODIFY:**
  - No protected files yet - this is the foundation task
  - After completion, other tasks must not modify core Go project structure

  **🤖 AI EXACT IMPLEMENTATION GUIDE:**

  **STEP 1: Initialize Go Project**
  ```bash
  go mod init lxc-terminal
  ```

  **STEP 2: Create Directory Structure**
  ```bash
  mkdir -p cmd/server
  mkdir -p internal/{lxc,pty,tmux,websocket}
  mkdir -p web/{css,js,assets/icons}
  mkdir -p tests/{unit,integration}
  mkdir -p docs
  ```

  **STEP 3: Create Main Server File (cmd/server/main.go)**
  ```go
  package main

  import (
      "fmt"
      "log"
      "net/http"
  )

  func main() {
      // Serve static files from web directory
      fs := http.FileServer(http.Dir("./web/"))
      http.Handle("/", fs)

      // WebSocket endpoint (placeholder for T2.1)
      http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
          fmt.Fprintf(w, "WebSocket endpoint - not implemented yet (T2.1)")
      })

      fmt.Println("LXC Terminal PWA Server starting on :8080")
      log.Fatal(http.ListenAndServe(":8080", nil))
  }
  ```

  **STEP 4: Create Package Stubs**
  - internal/lxc/manager.go - LXC container interface
  - internal/tmux/manager.go - tmux session interface
  - internal/websocket/handler.go - WebSocket server
  - internal/pty/terminal.go - PTY handling

  **STEP 5: Create Basic Web Structure**
  - web/index.html - Basic HTML page
  - web/css/ - CSS directory (T1.3 will populate)
  - web/js/ - JavaScript directory (T2.3 will populate)

  **🚫 FORBIDDEN - SCOPE CREEP PREVENTION:**
  - ❌ NO Docker integration (we use LXC)
  - ❌ NO authentication system implementation
  - ❌ NO database integration or ORM setup
  - ❌ NO frontend JavaScript beyond basic structure
  - ❌ NO actual WebSocket implementation (T2.1's job)
  - ❌ NO tmux integration code (T2.4's job)
  - ❌ NO PWA features like manifest/service worker (T1.2's job)
  - ❌ NO CSS styling or themes (T1.3's job)
  - ❌ NO testing frameworks beyond basic structure
  - ❌ NO deployment configuration beyond basic Dockerfile
  - ❌ NO logging or monitoring systems
  - ❌ NO configuration management systems

  **✅ VERIFICATION CHECKLIST:**
  ```bash
  # 1. Go project compiles successfully
  go build ./cmd/server || echo "FAIL: Build broken"

  # 2. Directory structure exists correctly
  ls -la cmd/server/main.go internal/ web/ || echo "FAIL: Structure missing"

  # 3. HTTP server starts and serves basic page
  go run cmd/server/main.go &
  sleep 2
  curl -s http://localhost:8080 | grep -q "html\|HTTP" || echo "FAIL: Server not responding"
  pkill -f "go run"

  # 4. Git repository setup
  git init
  git add .
  git status || echo "FAIL: Git not working"

  # 5. Go modules work correctly
  go mod tidy
  go mod verify || echo "FAIL: Go modules broken"
  ```

  **🎯 SUCCESS CRITERIA:**
  - Complete Go project structure created and compiles ✓
  - HTTP server starts on :8080 and serves basic content ✓
  - All package directories and stub files created ✓
  - Git repository initialized and ready for commits ✓
  - No scope creep - only foundation structure created ✓

  **📝 HANDOFF SUMMARY TEMPLATE - UPDATE WHEN COMPLETE:**
  ```markdown
  ## 🤖 T1.1 COMPLETION SUMMARY

  **COMPLETED BY:** [AI Model Name]
  **DATE:** [Date]
  **STATUS:** [Generated/Integrated/Tested/Complete]

  **WHAT WAS BUILT:**
  - Complete Go project structure in cmd/, internal/, web/ directories
  - Basic HTTP server in cmd/server/main.go serving on :8080
  - Package stubs for LXC, tmux, WebSocket, PTY integration
  - Git repository initialized with proper .gitignore
  - Go modules configured with go.mod

  **SCOPE DISCIPLINE MAINTAINED:**
  ✅ NO scope creep beyond Go project foundation
  ✅ NO authentication or security features added
  ✅ NO actual functionality beyond HTTP server
  ✅ NO frontend features beyond basic structure

  **BUILD STATUS:**
  ✅ go build ./cmd/server compiles successfully
  ✅ HTTP server starts and responds on :8080
  ✅ All package directories created correctly
  ✅ Git repository initialized and ready

  **INTEGRATION STATUS:**
  ✅ Foundation ready for T1.2 PWA features
  ✅ Package stubs ready for T2.1-T2.4 implementation
  ✅ Web directory ready for T1.3 CSS framework

  **NEXT AI TASKS CAN NOW:**
  - T1.2: Add PWA manifest and service worker to web/
  - T1.4: Define WebSocket protocol in internal/websocket/
  - T2.1: Implement WebSocket handler using foundation structure

  **CRITICAL NOTES FOR NEXT DEVELOPER:**
  - Do not modify cmd/server/main.go structure - other tasks extend it
  - Package interfaces in internal/ are stubs - implement in respective tasks
  - Web directory is minimal - T1.2/T1.3 will populate it
  ```

  **🔄 CONTINUOUS INTEGRATION REQUIREMENTS:**
  - Go build must succeed before task completion
  - Basic HTTP server must start without errors
  - All directory structure must be in place
  - Git repository must be properly initialized
  - No compilation errors or warnings allowed

  **📱 MOBILE TESTING REQUIREMENTS:**
  - HTTP server accessible from mobile browser on local network
  - Basic page loads correctly on iOS Safari and Android Chrome
  - No immediate mobile compatibility issues in server setup
EOF

echo "✅ T1.1 fixed with proper YAML indentation and 150+ lines of context!"