package websocket

// Message types for WebSocket protocol
const (
	// Terminal I/O
	MessageTypeTerminalInput  = "terminal_input"
	MessageTypeTerminalOutput = "terminal_output"
	MessageTypeTerminalResize = "terminal_resize"

	// Connection status
	MessageTypeConnect    = "connect"
	MessageTypeDisconnect = "disconnect"
	MessageTypeStatus     = "status"

	// LXC container management
	MessageTypeContainerList   = "container_list"
	MessageTypeContainerSelect = "container_select"
	MessageTypeContainerInfo   = "container_info"

	// Error handling
	MessageTypeError = "error"
)

// TerminalInput represents input data sent to terminal
type TerminalInput struct {
	Data string `json:"data"`
}

// TerminalOutput represents output data from terminal
type TerminalOutput struct {
	Data string `json:"data"`
}

// TerminalResize represents terminal resize command
type TerminalResize struct {
	Width  int `json:"width"`
	Height int `json:"height"`
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

// ContainerSelectMessage represents container selection request
type ContainerSelectMessage struct {
	ContainerName string `json:"container_name"`
}

// ContainerInfoMessage represents container information request
type ContainerInfoMessage struct {
	ContainerName string `json:"container_name"`
}