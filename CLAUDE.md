# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**Dash Terminal** is an LXC Terminal PWA with native tmux integration - a Go PTY server + xterm.js frontend designed as a premium web-based terminal for LXC containers.

## Architecture

- **Backend:** Go + WebSocket + PTY
- **Frontend:** Vanilla JS + xterm.js + Web Components
- **Session Management:** Native tmux integration (sessions → windows → panes)
- **Container Integration:** LXC exec
- **Deployment:** PWA + Service Worker on Tailscale HTTPS

## Development Commands

```bash
# Start development server
go run cmd/server/main.go

# Build application
go build -o dash-terminal cmd/server/main.go

# Run tests
go test ./...

# Lint code
gofmt -l . && go vet ./...
```

## DPPM Project Management

This project uses DPPM for task management. **CRITICAL:** Always check DPPM status first:

```bash
dppm status project dash-terminal
```

### Task Execution Order
Follow this **exact** order - AI instances MUST NOT skip or reorder:

**Phase 1 - Foundation:**
1. T1.1 project-setup (Go structure + tests)
2. T1.2 pwa-foundation (PWA manifest + service worker)
3. T1.3 tmux-ui-framework (CSS tabs only)
4. T1.4 tmux-websocket-protocol (Go message structs)

**Phase 2 - Core Functionality:**
5. T2.1 websocket-handler (Go WebSocket server)
6. T2.2 lxc-manager (Go LXC interface)
7. T2.3 terminal-wrapper (xterm.js integration)
8. T2.4 native-tmux-manager (Go tmux control)

**Phase 3 - Mobile Polish:**
9. T3.1 session-tabs (JavaScript components)
10. T3.2 window-tabs (JavaScript components)
11. T3.3 pane-grid-system (Layout rendering)
12. T3.4 touch-gestures (Mobile interactions)

## Target File Structure

```
cmd/server/main.go          # Go HTTP server entry point
internal/
  lxc/                      # LXC container management
  pty/                      # PTY session handling
  tmux/                     # Native tmux integration
  websocket/                # WebSocket message handling
web/
  index.html               # PWA application
  manifest.json            # PWA manifest
  sw.js                    # Service worker
  css/tmux-ui.css          # Native tmux UI styling
  js/                      # xterm.js + WebSocket client
```

## Tailscale HTTPS Setup

**Domain:** developer-dash.tail2d448.ts.net
**Base URL:** https://developer-dash.tail2d448.ts.net
**WebSocket:** wss://developer-dash.tail2d448.ts.net/ws

All testing and verification must use the Tailscale HTTPS endpoints.

## Dependencies

- Go >= 1.21
- LXC >= 5.0
- tmux >= 3.0
- github.com/gorilla/websocket v1.5.0
- github.com/creack/pty v1.1.18

## Development Rules

- **Test-Driven Development** required for all components
- **Anti-scope creep** discipline - strict MVP enforcement
- **NEVER** start with PWA/frontend - backend foundation first
- Follow the three-level tmux UI system: Sessions → Windows → Panes
- Performance targets: Sub-50ms terminal response, <200ms session switching, 60fps animations, <2MB PWA size

## Critical Mobile Development Rules

- **Touch + Click Unity**: NEVER code touch and click separately - use unified event handlers
- **PWA Cache Versioning**: ALWAYS bump cache version in sw.js when updating JavaScript/CSS
- **Mobile Testing**: PWA updates require cache invalidation or reinstall to see changes
- **Responsive Safe Areas**: Use `100dvh` and `max(env(safe-area-inset-top), 8px)` for mobile viewport

## MVP Scope Boundaries

**Included:**
- Basic terminal with xterm.js
- LXC container selection
- Native tmux session/window/pane management
- PWA installation
- Mobile-responsive design

**Explicitly Forbidden:**
- Authentication systems
- File upload/download
- Multi-server support
- Theming system
- Advanced tmux features beyond basic sessions