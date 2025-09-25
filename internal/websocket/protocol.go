package websocket

// Message types for tmux WebSocket protocol
const (
	// Session management
	MessageTypeSessionCreate = "session_create"
	MessageTypeSessionList   = "session_list"
	MessageTypeSessionSelect = "session_select"

	// Window management within sessions
	MessageTypeWindowCreate = "window_create"
	MessageTypeWindowList   = "window_list"
	MessageTypeWindowSelect = "window_select"

	// Pane management within windows
	MessageTypePaneCreate = "pane_create"
	MessageTypePaneList   = "pane_list"
	MessageTypePaneSelect = "pane_select"

	// Terminal I/O
	MessageTypeTerminalInput  = "terminal_input"
	MessageTypeTerminalOutput = "terminal_output"
	MessageTypeTerminalResize = "terminal_resize"

	// Connection status
	MessageTypeConnect    = "connect"
	MessageTypeDisconnect = "disconnect"
	MessageTypeStatus     = "status"

	// Error handling
	MessageTypeError = "error"
)

// Session represents a tmux session
type Session struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Window represents a tmux window within a session
type Window struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	SessionID string `json:"session_id"`
}

// Pane represents a tmux pane within a window
type Pane struct {
	ID       string `json:"id"`
	WindowID string `json:"window_id"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
}

// TerminalInput represents input data sent to terminal
type TerminalInput struct {
	PaneID string `json:"pane_id"`
	Data   string `json:"data"`
}

// TerminalOutput represents output data from terminal
type TerminalOutput struct {
	PaneID string `json:"pane_id"`
	Data   string `json:"data"`
}

// TerminalResize represents terminal resize command
type TerminalResize struct {
	PaneID string `json:"pane_id"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

// ErrorMessage represents an error response
type ErrorMessage struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// StatusMessage represents connection status
type StatusMessage struct {
	Connected bool   `json:"connected"`
	Message   string `json:"message"`
}
