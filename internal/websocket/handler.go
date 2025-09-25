package websocket

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"lxc-terminal/internal/lxc"
	"lxc-terminal/internal/pty"
	"lxc-terminal/internal/tmux"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins for development - restrict in production
		return true
	},
}

// Handler manages WebSocket connections
type Handler struct {
	connections       map[*websocket.Conn]bool
	lxcManager        *lxc.Manager
	ptyTerminal       *pty.Terminal
	tmuxManagers      map[string]*tmux.Manager // containerName -> tmux manager
	activeSessions    map[*websocket.Conn]string // conn -> sessionID
	mutex             sync.RWMutex
}

// NewHandler creates a new WebSocket handler
func NewHandler() *Handler {
	return &Handler{
		connections:    make(map[*websocket.Conn]bool),
		lxcManager:     lxc.NewManager(),
		ptyTerminal:    pty.NewTerminal(),
		tmuxManagers:   make(map[string]*tmux.Manager),
		activeSessions: make(map[*websocket.Conn]string),
	}
}

// HandleWebSocket handles WebSocket connection requests
func (h *Handler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}
	defer conn.Close()

	// Register connection
	h.connections[conn] = true
	defer func() {
		delete(h.connections, conn)

		// Clean up active session for this connection
		h.mutex.Lock()
		if sessionID, exists := h.activeSessions[conn]; exists {
			if session, ok := h.ptyTerminal.GetSession(sessionID); ok {
				session.Close()
			}
			delete(h.activeSessions, conn)
		}
		h.mutex.Unlock()
	}()

	// Send connection status
	statusMsg := Message{
		Type: MessageTypeStatus,
		Payload: StatusMessage{
			Connected: true,
			Message:   "Connected to Dash Terminal",
		},
	}

	if err := conn.WriteJSON(statusMsg); err != nil {
		log.Printf("Error sending status message: %v", err)
		return
	}

	// Handle incoming messages
	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// Process message based on type
		h.handleMessage(conn, &msg)
	}
}

// handleMessage processes incoming WebSocket messages
func (h *Handler) handleMessage(conn *websocket.Conn, msg *Message) {
	switch msg.Type {
	case MessageTypeConnect:
		h.handleConnect(conn)
	case MessageTypeSessionList:
		h.handleSessionList(conn)
	case MessageTypeSessionCreate:
		h.handleSessionCreate(conn, msg)
	case MessageTypeSessionSelect:
		h.handleSessionSelect(conn, msg)
	case MessageTypeWindowList:
		h.handleWindowList(conn, msg)
	case MessageTypeWindowCreate:
		h.handleWindowCreate(conn, msg)
	case MessageTypeWindowSelect:
		h.handleWindowSelect(conn, msg)
	case MessageTypePaneList:
		h.handlePaneList(conn, msg)
	case MessageTypePaneCreate:
		h.handlePaneCreate(conn, msg)
	case MessageTypePaneSelect:
		h.handlePaneSelect(conn, msg)
	case MessageTypeTerminalInput:
		h.handleTerminalInput(conn, msg)
	case MessageTypeTerminalResize:
		h.handleTerminalResize(conn, msg)
	case MessageTypeContainerList:
		h.handleContainerList(conn)
	case MessageTypeContainerSelect:
		h.handleContainerSelect(conn, msg)
	case MessageTypeContainerInfo:
		h.handleContainerInfo(conn, msg)
	default:
		h.sendError(conn, 400, "Unknown message type: "+msg.Type)
	}
}

// handleConnect handles connection requests
func (h *Handler) handleConnect(conn *websocket.Conn) {
	response := Message{
		Type: MessageTypeStatus,
		Payload: StatusMessage{
			Connected: true,
			Message:   "Successfully connected",
		},
	}
	conn.WriteJSON(response)
}

// Placeholder handlers for T2.2+ implementation
func (h *Handler) handleSessionList(conn *websocket.Conn) {
	// TODO: T2.4 will implement actual tmux session listing
	response := Message{
		Type:    MessageTypeSessionList,
		Payload: []Session{},
	}
	conn.WriteJSON(response)
}

func (h *Handler) handleSessionCreate(conn *websocket.Conn, msg *Message) {
	// TODO: T2.4 will implement session creation
	h.sendError(conn, 501, "Session creation not yet implemented")
}

func (h *Handler) handleSessionSelect(conn *websocket.Conn, msg *Message) {
	// TODO: T2.4 will implement session selection
	h.sendError(conn, 501, "Session selection not yet implemented")
}

func (h *Handler) handleWindowList(conn *websocket.Conn, msg *Message) {
	// TODO: T2.4 will implement window listing
	response := Message{
		Type:    MessageTypeWindowList,
		Payload: []Window{},
	}
	conn.WriteJSON(response)
}

func (h *Handler) handleWindowCreate(conn *websocket.Conn, msg *Message) {
	// TODO: T2.4 will implement window creation
	h.sendError(conn, 501, "Window creation not yet implemented")
}

func (h *Handler) handleWindowSelect(conn *websocket.Conn, msg *Message) {
	// TODO: T2.4 will implement window selection
	h.sendError(conn, 501, "Window selection not yet implemented")
}

func (h *Handler) handlePaneList(conn *websocket.Conn, msg *Message) {
	// TODO: T2.4 will implement pane listing
	response := Message{
		Type:    MessageTypePaneList,
		Payload: []Pane{},
	}
	conn.WriteJSON(response)
}

func (h *Handler) handlePaneCreate(conn *websocket.Conn, msg *Message) {
	// TODO: T2.4 will implement pane creation
	h.sendError(conn, 501, "Pane creation not yet implemented")
}

func (h *Handler) handlePaneSelect(conn *websocket.Conn, msg *Message) {
	// TODO: T2.4 will implement pane selection
	h.sendError(conn, 501, "Pane selection not yet implemented")
}

func (h *Handler) handleTerminalInput(conn *websocket.Conn, msg *Message) {
	// TODO: T2.3 will implement terminal input handling
	h.sendError(conn, 501, "Terminal input not yet implemented")
}

func (h *Handler) handleTerminalResize(conn *websocket.Conn, msg *Message) {
	// TODO: T2.3 will implement terminal resize handling
	h.sendError(conn, 501, "Terminal resize not yet implemented")
}

// LXC container handlers

// handleContainerList handles container list requests
func (h *Handler) handleContainerList(conn *websocket.Conn) {
	containers, err := h.lxcManager.ListContainers()
	if err != nil {
		h.sendError(conn, 500, "Failed to list containers: "+err.Error())
		return
	}

	response := Message{
		Type:    MessageTypeContainerList,
		Payload: containers,
	}
	conn.WriteJSON(response)
}

// handleContainerSelect handles container selection requests
func (h *Handler) handleContainerSelect(conn *websocket.Conn, msg *Message) {
	var selectMsg ContainerSelectMessage
	if err := h.parseMessagePayload(msg, &selectMsg); err != nil {
		h.sendError(conn, 400, "Invalid container select message: "+err.Error())
		return
	}

	// Check if container exists and is running
	running, err := h.lxcManager.IsContainerRunning(selectMsg.ContainerName)
	if err != nil {
		h.sendError(conn, 404, "Container not found: "+selectMsg.ContainerName)
		return
	}

	if !running {
		h.sendError(conn, 400, "Container is not running: "+selectMsg.ContainerName)
		return
	}

	// For T2.2, we just confirm selection is valid
	// Actual connection will be handled in T2.3/T2.4
	response := Message{
		Type: MessageTypeStatus,
		Payload: StatusMessage{
			Connected: true,
			Message:   "Container " + selectMsg.ContainerName + " selected successfully",
		},
	}
	conn.WriteJSON(response)
}

// handleContainerInfo handles container info requests
func (h *Handler) handleContainerInfo(conn *websocket.Conn, msg *Message) {
	var infoMsg ContainerInfoMessage
	if err := h.parseMessagePayload(msg, &infoMsg); err != nil {
		h.sendError(conn, 400, "Invalid container info message: "+err.Error())
		return
	}

	container, err := h.lxcManager.GetContainer(infoMsg.ContainerName)
	if err != nil {
		h.sendError(conn, 404, "Container not found: "+infoMsg.ContainerName)
		return
	}

	response := Message{
		Type:    MessageTypeContainerInfo,
		Payload: container,
	}
	conn.WriteJSON(response)
}

// parseMessagePayload parses message payload into target struct
func (h *Handler) parseMessagePayload(msg *Message, target interface{}) error {
	// Convert payload to JSON and back to parse into struct
	// This handles the interface{} -> specific struct conversion
	payloadBytes, err := json.Marshal(msg.Payload)
	if err != nil {
		return err
	}
	return json.Unmarshal(payloadBytes, target)
}

// sendError sends an error message to the client
func (h *Handler) sendError(conn *websocket.Conn, code int, message string) {
	errorMsg := Message{
		Type: MessageTypeError,
		Payload: ErrorMessage{
			Code:    code,
			Message: message,
		},
	}
	conn.WriteJSON(errorMsg)
}

// Message represents a WebSocket message structure
type Message struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}
