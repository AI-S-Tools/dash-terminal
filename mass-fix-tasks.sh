#!/bin/bash

# Mass fix all 11 tasks with anti-scope creep system
# Based on T1.1 and T1.4 successful templates

DPPM_BASE="/home/ubuntu/Dropbox/project-management/projects/dash-terminal"

echo "🔧 Mass updating all tasks with anti-scope creep system..."

# Function to create task description
create_task_description() {
    local task_id="$1"
    local title="$2"
    local scope="$3"
    local prerequisites="$4"
    local previous_summary="$5"
    local file_mode="$6"
    local protected_files="$7"
    local exact_implementation="$8"
    local forbidden="$9"
    local verification="${10}"
    local success_criteria="${11}"

    cat > temp_description.md << EOF
## $title - $scope

**⚠️ STOP! CHECK DEPENDENCIES FIRST:**
\`\`\`bash
$prerequisites
\`\`\`

**📋 SUMMARY FROM PREVIOUS TASKS:**
$previous_summary

**📂 FILE STATUS CHECK - RUN THIS FIRST:**
\`\`\`bash
# CHECK IF THIS TASK IS ALREADY DONE
# [Task-specific existence checks will be added]
\`\`\`

**$file_mode**

**🔒 PROTECTED FILES - DO NOT MODIFY:**
$protected_files

**🤖 AI EXACT IMPLEMENTATION:**
$exact_implementation

**🚫 FORBIDDEN - DO NOT ADD:**
$forbidden

**✅ VERIFICATION CHECKLIST:**
\`\`\`bash
$verification
\`\`\`

**🎯 SUCCESS CRITERIA:**
$success_criteria

**📝 HANDOFF SUMMARY - UPDATE WHEN COMPLETE:**
\`\`\`
## 🤖 $task_id COMPLETION SUMMARY

**COMPLETED BY:** [AI Model Name]
**DATE:** [Date]
**STATUS:** Ready for next task

**WHAT WAS BUILT:**
- [List specific deliverables]

**SCOPE DISCIPLINE MAINTAINED:**
✅ NO scope creep beyond defined boundaries
✅ NO features from other tasks implemented

**BUILD STATUS:**
✅ go build ./cmd/server compiles
✅ go test ./... passes

**NEXT AI TASK:**
[Next task can now proceed with these foundations]
\`\`\`
EOF
}

# T1.2: PWA Foundation
echo "Fixing T1.2..."
cat > "$DPPM_BASE/phases/p1-foundation/tasks/t1-2-pwa-foundation.yaml" << 'EOF'
id: t1-2-pwa-foundation
title: 'T1.2: PWA Foundation'
project_id: dash-terminal
phase_id: p1-foundation
status: todo
priority: high
reporter: dppm-user
created: "2025-09-25"
updated: "2025-09-25"
dependency_ids: ["t1-1-project-setup"]
description: |
  ## T1.2: PWA Foundation - MANIFEST & SERVICE WORKER ONLY

  **⚠️ STOP! CHECK DEPENDENCIES FIRST:**
  ```bash
  # YOU CANNOT START UNTIL T1.1 IS DONE:
  dppm task show t1-1-project-setup --project dash-terminal | grep "status: done" || { echo "ERROR: T1.1 not done"; exit 1; }
  ls -la cmd/server/main.go || { echo "ERROR: T1.1 incomplete - no Go server"; exit 1; }
  ls -la web/ || { echo "ERROR: T1.1 incomplete - no web directory"; exit 1; }
  ```

  **📋 SUMMARY FROM PREVIOUS TASKS:**
  - ✅ **T1.1**: Go project structure with basic HTTP server on :8080 (serves "Not Implemented")
  - 📍 **You are here**: T1.2 - Add PWA installability to existing structure

  **📂 FILE STATUS CHECK - RUN THIS FIRST:**
  ```bash
  # CHECK IF ALREADY DONE
  ls -la web/manifest.json 2>/dev/null && { echo "WARNING: manifest.json EXISTS - task may be done"; exit 1; }
  ls -la web/sw.js 2>/dev/null && { echo "WARNING: service worker EXISTS - task may be done"; exit 1; }
  ```

  **🔄 FILE OPERATION MODE: CREATE NEW**
  - CREATE: web/manifest.json, web/sw.js, web/index.html, web/assets/icons
  - Use 'Write' tool for all new files
  - Never use 'Edit' tool - all files are new

  **🔒 PROTECTED FILES - DO NOT MODIFY:**
  - cmd/server/main.go (T1.1's responsibility)
  - go.mod (only T1.1 should modify)
  - internal/ directory (backend - not this task)
  - Any existing .go files

  **🤖 AI EXACT IMPLEMENTATION:**
  [Full PWA implementation from T1.2 template above]

  **🚫 FORBIDDEN - DO NOT ADD:**
  - ❌ NO xterm.js yet (T2.3's job)
  - ❌ NO WebSocket client code (T2.3's job)
  - ❌ NO terminal UI elements (T1.3's job)
  - ❌ NO JavaScript beyond service worker registration
  - ❌ NO CSS files (T1.3's job)
  - ❌ NO backend modifications
  - ❌ NO npm/package.json
  - ❌ NO build tools
  - ❌ NO authentication
  - ❌ NO database setup

  **✅ VERIFICATION CHECKLIST:**
  ```bash
  # 1. PWA files exist
  ls -la web/manifest.json web/sw.js web/index.html || echo "FAIL"

  # 2. PWA installability (manual test in Chrome)
  # Open http://localhost:8080 in Chrome DevTools > Application > Manifest
  # Should show "Installable" status

  # 3. Service worker registration
  curl -s http://localhost:8080 | grep "serviceWorker.register" || echo "FAIL"
  ```

  **🎯 SUCCESS CRITERIA:**
  - PWA installs on mobile devices ✓
  - Service worker caches for offline use ✓
  - Shows "LXC Terminal PWA - Ready for Installation" ✓
  - NO terminal functionality yet ✓

EOF

# T1.3: Native Tmux UI Framework
echo "Fixing T1.3..."
cat > "$DPPM_BASE/phases/p1-foundation/tasks/t1-3-tmux-ui-framework.yaml" << 'EOF'
id: t1-3-tmux-ui-framework
title: 'T1.3: Native Tmux UI Framework'
project_id: dash-terminal
phase_id: p1-foundation
status: todo
priority: high
reporter: dppm-user
created: "2025-09-25"
updated: "2025-09-25"
dependency_ids: ["t1-2-pwa-foundation"]
description: |
  ## T1.3: Native Tmux UI Framework - CSS TAB STRUCTURE ONLY

  **⚠️ STOP! CHECK DEPENDENCIES FIRST:**
  ```bash
  # YOU CANNOT START UNTIL T1.2 IS DONE:
  dppm task show t1-2-pwa-foundation --project dash-terminal | grep "status: done" || { echo "ERROR: T1.2 not done"; exit 1; }
  ls -la web/index.html || { echo "ERROR: T1.2 incomplete - no HTML"; exit 1; }
  ls -la web/manifest.json || { echo "ERROR: T1.2 incomplete - no PWA"; exit 1; }
  ```

  **📋 SUMMARY FROM PREVIOUS TASKS:**
  - ✅ **T1.1**: Go HTTP server running on :8080
  - ✅ **T1.2**: PWA with index.html showing "Ready for Installation"
  - 📍 **You are here**: T1.3 - Add CSS tab structure (no JavaScript yet)

  **📂 FILE STATUS CHECK - RUN THIS FIRST:**
  ```bash
  # CHECK IF ALREADY DONE
  ls -la web/css/tmux-ui.css 2>/dev/null && { echo "WARNING: CSS EXISTS - task may be done"; exit 1; }
  grep "session-tabs" web/index.html 2>/dev/null && { echo "WARNING: Tabs already in HTML"; exit 1; }
  ```

  **🔄 FILE OPERATION MODE: EXTEND EXISTING + CREATE NEW**
  - EXTEND: Modify web/index.html body content
  - CREATE: New web/css/tmux-ui.css file
  - Use 'Read' on index.html first, then 'Edit'
  - Use 'Write' for CSS file

  **🔒 PROTECTED FILES - DO NOT MODIFY:**
  - web/sw.js (service worker - T1.2's work)
  - web/manifest.json (PWA config - T1.2's work)
  - cmd/server/main.go (backend - not this task)
  - internal/ directory (backend - not this task)

  **🤖 AI EXACT IMPLEMENTATION:**
  [Full CSS framework from template above - session tabs, window tabs, pane grid]

  **🚫 FORBIDDEN - DO NOT ADD:**
  - ❌ NO JavaScript event handlers (T3.1-T3.4's job)
  - ❌ NO WebSocket connections (T2.1's job)
  - ❌ NO xterm.js integration (T2.3's job)
  - ❌ NO touch gestures (T3.4's job)
  - ❌ NO actual terminal content (T2.3's job)
  - ❌ NO animations or transitions yet
  - ❌ NO responsive breakpoints beyond mobile-first
  - ❌ NO theme switching
  - ❌ NO extra CSS files

  **✅ VERIFICATION CHECKLIST:**
  ```bash
  # 1. CSS file created
  ls -la web/css/tmux-ui.css || echo "FAIL: No CSS"

  # 2. HTML updated with tabs
  grep "session-tabs" web/index.html || echo "FAIL: No tabs in HTML"

  # 3. Visual test
  # Open http://localhost:8080 - should see 3-level tab structure
  ```

  **🎯 SUCCESS CRITERIA:**
  - Three-level tab structure visible ✓
  - Session tabs, window tabs, pane grid ✓
  - Mobile-responsive layout ✓
  - NO functionality yet (just visual) ✓

EOF

echo "Created T1.2 and T1.3 fixes. Continue with remaining 9 tasks..."
echo "Files updated:"
ls -la "$DPPM_BASE/phases/p1-foundation/tasks/"