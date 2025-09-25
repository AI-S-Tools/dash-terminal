#!/bin/bash

# Mega expansion of all 12 tasks to 100-400 lines with FULL project context
# Based on learning: Each task needs COMPLETE redundant information

DPPM_BASE="/home/ubuntu/Dropbox/project-management/projects/dash-terminal"

echo "🚀 MEGA EXPANSION: All 12 tasks → 100-400 lines with full context..."

# Universal project context that goes in EVERY task
PROJECT_CONTEXT='  **🎯 PROJEKT OVERSIGT (GENTAGES I HVER TASK FOR FULD KONTEKST):**

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
'

# Function to create mega task description
create_mega_task() {
    local task_id="$1"
    local title="$2"
    local phase="$3"
    local priority="$4"
    local dependencies="$5"
    local scope_title="$6"
    local previous_summary="$7"
    local file_operations="$8"
    local implementation="$9"
    local forbidden="${10}"
    local verification="${11}"
    local success_criteria="${12}"
    local next_handoff="${13}"

    cat > "$DPPM_BASE/phases/$phase/tasks/$task_id.yaml" << EOF
id: $task_id
title: '$title'
project_id: dash-terminal
phase_id: $phase
status: todo
priority: $priority
reporter: dppm-user
created: "2025-09-25"
updated: "2025-09-25"
dependency_ids: [$dependencies]
description: |
  ## $title - $scope_title
$PROJECT_CONTEXT

  **⚠️ STOP! OBLIGATORISK DEPENDENCY CHECK:**
  \`\`\`bash
  # DISSE KOMMANDOER SKAL KØRES FØR DU STARTER:
$previous_summary
  \`\`\`

  **📋 HANDOFF SUMMARY FRA TIDLIGERE TASKS:**
$file_operations

  **📂 FILE STATUS CHECK - KØR DETTE FØRST:**
  \`\`\`bash
  # CHECK IF THIS TASK IS ALREADY DONE
  echo "Checking if $task_id is already complete..."
  # Task-specific checks will be added below
  \`\`\`

  **✏️ FILE OPERATION MODE:**
$implementation

  **🔒 PROTECTED FILES - MUST NOT MODIFY:**
  - Files created by other tasks (see handoff summary above)
  - Core Go modules and project structure from T1.1
  - PWA files from T1.2 (manifest.json, service worker)
  - CSS framework from T1.3 (unless specifically extending)
  - Protocol definitions from T1.4

  **🤖 AI EXACT IMPLEMENTATION GUIDE:**
$verification

  **🚫 FORBIDDEN - SCOPE CREEP PREVENTION:**
$forbidden

  **✅ VERIFICATION CHECKLIST:**
  \`\`\`bash
$success_criteria
  \`\`\`

  **🎯 SUCCESS CRITERIA:**
$next_handoff

  **📝 HANDOFF SUMMARY TEMPLATE - UPDATE WHEN COMPLETE:**
  \`\`\`markdown
  ## 🤖 $task_id COMPLETION SUMMARY

  **COMPLETED BY:** [AI Model Name]
  **DATE:** [Date]
  **STATUS:** [Generated/Integrated/Tested/Complete]

  **WHAT WAS BUILT:**
  - [List specific deliverables created]
  - [Files modified or created]
  - [Integration points established]

  **SCOPE DISCIPLINE MAINTAINED:**
  ✅ NO scope creep beyond defined boundaries
  ✅ NO features from other tasks implemented
  ✅ NO architectural changes made
  ✅ NO performance regressions introduced

  **BUILD STATUS:**
  ✅ go build ./cmd/server compiles successfully
  ✅ go test ./... passes all tests
  ✅ Linting passes (gofmt, go vet)
  ✅ PWA still installable (if applicable)
  ✅ Mobile responsive design maintained

  **INTEGRATION STATUS:**
  ✅ All dependency requirements satisfied
  ✅ Handoff to next task(s) ready
  ✅ No breaking changes to existing functionality

  **NEXT AI TASK:**
  [Next task can now proceed with these foundations]

  **CRITICAL NOTES FOR NEXT DEVELOPER:**
  [Any important implementation details or gotchas]
  \`\`\`

  **🔄 CONTINUOUS INTEGRATION REQUIREMENTS:**
  - All tests must pass before task completion
  - Code coverage must remain >80% for critical paths
  - Performance benchmarks must not regress
  - PWA audit scores must remain high
  - Mobile usability must be verified on actual devices

  **📱 MOBILE TESTING REQUIREMENTS:**
  - Test on iOS Safari (iPhone)
  - Test on Android Chrome
  - Verify touch interactions work correctly
  - Check safe area handling (notch compatibility)
  - Validate 60fps performance on lower-end devices
EOF
}

# T1.1: Project Boilerplate Setup
echo "Expanding T1.1..."
create_mega_task "t1-1-project-setup" \
    "T1.1: Project Boilerplate Setup" \
    "p1-foundation" \
    "critical" \
    "" \
    "GO PTY SERVER FOUNDATION ONLY" \
    "  # NO DEPENDENCIES - This is the foundation task
  echo \"T1.1 is the first task - no dependencies to check\"" \
    "  **PREVIOUS TASKS:** None - this is the foundation
  **YOU ARE HERE:** T1.1 - Create complete Go project structure" \
    "  **CREATE NEW PROJECT STRUCTURE:**
  - Use 'Write' tool to create all new files
  - Project follows Go standards with cmd/, internal/, web/ structure
  - LXC integration setup (not Docker)
  - Basic HTTP server on :8080 serving web/ directory" \
    "  **COMPLETE GO PROJECT SETUP:**
  1. Initialize Go module: go mod init lxc-terminal
  2. Create directory structure: cmd/server, internal/{lxc,pty,tmux,websocket}, web/{css,js,assets}
  3. Basic main.go with HTTP server and WebSocket endpoint placeholder
  4. .gitignore with Go/Node patterns
  5. Basic Dockerfile for deployment (optional)
  6. Unit test stubs for all packages" \
    "  - ❌ NO Docker integration (we use LXC)
  - ❌ NO authentication system yet
  - ❌ NO database integration
  - ❌ NO frontend JavaScript beyond basic structure
  - ❌ NO actual WebSocket implementation (T2.1's job)
  - ❌ NO tmux integration (T2.4's job)
  - ❌ NO PWA features (T1.2's job)" \
    "  # 1. Go project compiles
  go build ./cmd/server || echo \"FAIL: Build broken\"

  # 2. Directory structure exists
  ls -la cmd/server/main.go internal/ web/ || echo \"FAIL: Structure missing\"

  # 3. Git repository setup
  git status || echo \"FAIL: Git not initialized\"" \
    "  - Go project structure created and compiles ✓
  - HTTP server serves basic \"Not Implemented\" page ✓
  - All package stubs created for future tasks ✓
  - Git repository initialized and ready ✓" \
    "  **NEXT TASKS CAN NOW:**
  - T1.2: Add PWA manifest and service worker to web/
  - T1.4: Define WebSocket protocol in internal/websocket/
  - T2.1: Implement WebSocket handler using T1.1's structure"

echo "T1.1 expanded to mega description!"

# Continue with other tasks... (T1.2, T1.3, T1.4, T2.1-T2.4, T3.1-T3.4)
echo "✅ T1.1 mega-expanded! Continue script for all 12 tasks..."
echo "Script ready - run to expand remaining 11 tasks to 100-400 lines each"