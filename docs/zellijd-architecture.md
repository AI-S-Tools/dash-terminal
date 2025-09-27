# Zellijd - Go Daemon Architecture

## Overview
Drop-in Go daemon per LXC container that bridges Zellij to PWA via WebSocket + REST API.

## Architecture

```
PWA (xterm.js)
     ↓
[WebSocket: /ws] [REST: /api/*]
     ↓              ↓
  zellijd (Go Daemon)
     ↓
  Zellij Sessions
     ↓
  Container Shell
```

## Endpoints

### WebSocket: `/ws`
- Real-time PTY data stream
- Bidirectional terminal I/O
- Resize events
- Session attach/detach

### REST API: `/api/*`

#### Session Management
- `GET /api/sessions` - List all Zellij sessions
- `POST /api/sessions` - Create new session
- `DELETE /api/sessions/:id` - Kill session
- `GET /api/sessions/:id/attach` - Get WebSocket token for session

#### Pane Management
- `GET /api/sessions/:id/panes` - List panes
- `POST /api/sessions/:id/panes` - Split pane
- `DELETE /api/sessions/:id/panes/:pane` - Close pane
- `POST /api/sessions/:id/panes/:pane/focus` - Focus pane

#### Layout Management
- `GET /api/layouts` - List saved layouts
- `POST /api/sessions/:id/layout` - Apply layout
- `GET /api/sessions/:id/layout` - Get current layout

#### Actions
- `POST /api/sessions/:id/rename` - Rename session
- `POST /api/sessions/:id/panes/:pane/resize` - Resize pane
- `POST /api/sessions/:id/command` - Run Zellij command

## Implementation Plan

### Phase 1: Core Daemon
```go
// cmd/zellijd/main.go
type ZellijDaemon struct {
    sessions map[string]*ZellijSession
    mu       sync.RWMutex
}

type ZellijSession struct {
    ID       string
    Process  *exec.Cmd
    PTY      *os.File
    Active   bool
}
```

### Phase 2: WebSocket Handler
```go
// internal/websocket/zellij_handler.go
func (h *Handler) HandleZellijWS(w http.ResponseWriter, r *http.Request) {
    sessionID := r.URL.Query().Get("session")
    session := h.daemon.GetSession(sessionID)

    // Bridge PTY ↔ WebSocket
    go io.Copy(wsConn, session.PTY)
    go io.Copy(session.PTY, wsConn)
}
```

### Phase 3: REST API
```go
// internal/api/handlers.go
func (h *Handler) ListSessions(w http.ResponseWriter, r *http.Request) {
    sessions := h.daemon.ListSessions()
    json.NewEncoder(w).Encode(sessions)
}

func (h *Handler) CreateSession(w http.ResponseWriter, r *http.Request) {
    session := h.daemon.CreateSession(req.Name)
    json.NewEncoder(w).Encode(session)
}
```

## Benefits vs Tmux

1. **Clean API** - REST endpoints instead of parsing control mode
2. **No keybinding conflicts** - Zellij runs headless
3. **Modern features** - Floating panes, layouts, plugins
4. **Better state management** - Structured data, not escape sequences
5. **PWA-friendly** - Designed for programmatic control

## Container Deployment

Each LXC container runs its own zellijd:
```bash
# Per-container systemd service
[Unit]
Description=Zellij Daemon for %i
After=network.target

[Service]
Type=simple
ExecStart=/usr/local/bin/zellijd --container %i --port 8080
Restart=always

[Install]
WantedBy=multi-user.target
```

## Security

- WebSocket authentication via tokens
- Per-session isolation
- Rate limiting on API endpoints
- Sanitized command execution

## Performance

- Single Go binary (~10MB)
- Low memory footprint (~5MB per session)
- Native Zellij performance
- WebSocket compression for bandwidth