package websocket

import (
	"encoding/json"
	"fmt"
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

// handleSessionList handles session list requests - lists tmux sessions
func (h *Handler) handleSessionList(conn *websocket.Conn) {
	// Get current container for this connection
	h.mutex.RLock()
	sessionID, exists := h.activeSessions[conn]
	h.mutex.RUnlock()

	var containerName string
	if exists {
		// Extract container name from existing session
		// For now, we'll use empty string for host or extract from sessionID if needed
		containerName = "" // Host terminal
	}

	// Get or create tmux manager for this container
	tmuxManager := h.getTmuxManager(containerName)

	// List tmux sessions
	sessions, err := tmuxManager.ListSessions()
	if err != nil {
		h.sendError(conn, 500, "Failed to list tmux sessions: "+err.Error())
		return
	}

	// Convert tmux.Session to websocket.Session format
	wsSessions := make([]Session, len(sessions))
	for i, session := range sessions {
		wsSessions[i] = Session{
			ID:   session.Name,
			Name: session.Name,
		}
	}

	response := Message{
		Type:    MessageTypeSessionList,
		Payload: wsSessions,
	}
	conn.WriteJSON(response)
}

func (h *Handler) handleSessionCreate(conn *websocket.Conn, msg *Message) {
	var sessionMsg map[string]interface{}
	if err := h.parseMessagePayload(msg, &sessionMsg); err != nil {
		h.sendError(conn, 400, "Invalid session create message: "+err.Error())
		return
	}

	sessionName, ok := sessionMsg["name"].(string)
	if !ok || sessionName == "" {
		h.sendError(conn, 400, "Session name is required")
		return
	}

	// Get current container
	var containerName string
	h.mutex.RLock()
	if sessionID, exists := h.activeSessions[conn]; exists {
		containerName = "" // Host terminal for now
		_ = sessionID
	}
	h.mutex.RUnlock()

	// Get tmux manager
	tmuxManager := h.getTmuxManager(containerName)

	// Create tmux session
	err := tmuxManager.CreateSession(sessionName)
	if err != nil {
		h.sendError(conn, 500, "Failed to create tmux session: "+err.Error())
		return
	}

	// Send success response
	response := Message{
		Type: MessageTypeStatus,
		Payload: StatusMessage{
			Connected: true,
			Message:   "Tmux session created: " + sessionName,
		},
	}
	conn.WriteJSON(response)
}

func (h *Handler) handleSessionSelect(conn *websocket.Conn, msg *Message) {
	var sessionMsg map[string]interface{}
	if err := h.parseMessagePayload(msg, &sessionMsg); err != nil {
		h.sendError(conn, 400, "Invalid session select message: "+err.Error())
		return
	}

	sessionName, ok := sessionMsg["name"].(string)
	if !ok || sessionName == "" {
		h.sendError(conn, 400, "Session name is required")
		return
	}

	// Get current container
	var containerName string
	h.mutex.RLock()
	if sessionID, exists := h.activeSessions[conn]; exists {
		containerName = "" // Host terminal for now
		_ = sessionID
	}
	h.mutex.RUnlock()

	// Get tmux manager
	tmuxManager := h.getTmuxManager(containerName)

	// Select tmux session
	err := tmuxManager.SelectSession(sessionName)
	if err != nil {
		h.sendError(conn, 500, "Failed to select tmux session: "+err.Error())
		return
	}

	// Send success response
	response := Message{
		Type: MessageTypeStatus,
		Payload: StatusMessage{
			Connected: true,
			Message:   "Tmux session selected: " + sessionName,
		},
	}
	conn.WriteJSON(response)
}

func (h *Handler) handleWindowList(conn *websocket.Conn, msg *Message) {
	var windowMsg map[string]interface{}
	if err := h.parseMessagePayload(msg, &windowMsg); err != nil {
		h.sendError(conn, 400, "Invalid window list message: "+err.Error())
		return
	}

	sessionName, ok := windowMsg["session_name"].(string)
	if !ok || sessionName == "" {
		h.sendError(conn, 400, "Session name is required")
		return
	}

	// Get current container
	var containerName string
	h.mutex.RLock()
	if sessionID, exists := h.activeSessions[conn]; exists {
		containerName = "" // Host terminal for now
		_ = sessionID
	}
	h.mutex.RUnlock()

	// Get tmux manager
	tmuxManager := h.getTmuxManager(containerName)

	// List windows for session
	sessions, err := tmuxManager.ListSessions()
	if err != nil {
		h.sendError(conn, 500, "Failed to list tmux sessions: "+err.Error())
		return
	}

	// Find the specific session
	var targetSession *tmux.Session
	for i, session := range sessions {
		if session.Name == sessionName {
			targetSession = &sessions[i]
			break
		}
	}

	if targetSession == nil {
		h.sendError(conn, 404, "Session not found: "+sessionName)
		return
	}

	// Convert tmux.Window to websocket.Window format
	wsWindows := make([]Window, len(targetSession.Windows))
	for i, window := range targetSession.Windows {
		wsWindows[i] = Window{
			ID:        fmt.Sprintf("%d", window.ID),
			Name:      window.Name,
			SessionID: sessionName,
		}
	}

	response := Message{
		Type:    MessageTypeWindowList,
		Payload: wsWindows,
	}
	conn.WriteJSON(response)
}

func (h *Handler) handleWindowCreate(conn *websocket.Conn, msg *Message) {
	var windowMsg map[string]interface{}
	if err := h.parseMessagePayload(msg, &windowMsg); err != nil {
		h.sendError(conn, 400, "Invalid window create message: "+err.Error())
		return
	}

	sessionName, ok := windowMsg["session_name"].(string)
	if !ok || sessionName == "" {
		h.sendError(conn, 400, "Session name is required")
		return
	}

	windowName, ok := windowMsg["window_name"].(string)
	if !ok || windowName == "" {
		h.sendError(conn, 400, "Window name is required")
		return
	}

	// Get current container
	var containerName string
	h.mutex.RLock()
	if sessionID, exists := h.activeSessions[conn]; exists {
		containerName = "" // Host terminal for now
		_ = sessionID
	}
	h.mutex.RUnlock()

	// Get tmux manager
	tmuxManager := h.getTmuxManager(containerName)

	// Create tmux window
	err := tmuxManager.CreateWindow(sessionName, windowName)
	if err != nil {
		h.sendError(conn, 500, "Failed to create tmux window: "+err.Error())
		return
	}

	// Send success response
	response := Message{
		Type: MessageTypeStatus,
		Payload: StatusMessage{
			Connected: true,
			Message:   fmt.Sprintf("Tmux window created: %s in session %s", windowName, sessionName),
		},
	}
	conn.WriteJSON(response)
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
	var input TerminalInput
	if err := h.parseMessagePayload(msg, &input); err != nil {
		h.sendError(conn, 400, "Invalid terminal input: "+err.Error())
		return
	}

	// Get active session for this connection
	h.mutex.RLock()
	sessionID, exists := h.activeSessions[conn]
	h.mutex.RUnlock()

	if !exists {
		h.sendError(conn, 400, "No active terminal session")
		return
	}

	// Get PTY session
	session, ok := h.ptyTerminal.GetSession(sessionID)
	if !ok {
		h.sendError(conn, 404, "Terminal session not found")
		return
	}

	// Write input to PTY
	_, err := session.Write([]byte(input.Data))
	if err != nil {
		h.sendError(conn, 500, "Failed to write to terminal: "+err.Error())
		return
	}
}

func (h *Handler) handleTerminalResize(conn *websocket.Conn, msg *Message) {
	var resize TerminalResize
	if err := h.parseMessagePayload(msg, &resize); err != nil {
		h.sendError(conn, 400, "Invalid terminal resize: "+err.Error())
		return
	}

	// Get active session for this connection
	h.mutex.RLock()
	sessionID, exists := h.activeSessions[conn]
	h.mutex.RUnlock()

	if !exists {
		h.sendError(conn, 400, "No active terminal session")
		return
	}

	// Get PTY session
	session, ok := h.ptyTerminal.GetSession(sessionID)
	if !ok {
		h.sendError(conn, 404, "Terminal session not found")
		return
	}

	// Resize PTY
	err := session.Resize(resize.Width, resize.Height)
	if err != nil {
		h.sendError(conn, 500, "Failed to resize terminal: "+err.Error())
		return
	}

	// Send success response
	response := Message{
		Type: MessageTypeStatus,
		Payload: StatusMessage{
			Connected: true,
			Message:   fmt.Sprintf("Terminal resized to %dx%d", resize.Width, resize.Height),
		},
	}
	conn.WriteJSON(response)
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

	// Start terminal session in selected container
	sessionID := fmt.Sprintf("conn_%p", conn)

	// Start PTY session with bash in the container
	session, err := h.ptyTerminal.StartSession(sessionID, selectMsg.ContainerName, "bash")
	if err != nil {
		h.sendError(conn, 500, "Failed to start terminal session: "+err.Error())
		return
	}

	// Register active session
	h.mutex.Lock()
	h.activeSessions[conn] = sessionID
	h.mutex.Unlock()

	// Start reading from PTY and sending to WebSocket
	go h.pipeTerminalOutput(conn, session)

	// Send success response
	response := Message{
		Type: MessageTypeStatus,
		Payload: StatusMessage{
			Connected: true,
			Message:   "Terminal session started in container " + selectMsg.ContainerName,
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

// pipeTerminalOutput continuously reads from PTY and sends to WebSocket
func (h *Handler) pipeTerminalOutput(conn *websocket.Conn, session *pty.Session) {
	buffer := make([]byte, 1024)

	for session.IsRunning() {
		n, err := session.Read(buffer)
		if err != nil {
			if err != io.EOF {
				log.Printf("Error reading from PTY: %v", err)
			}
			break
		}

		if n > 0 {
			// Send terminal output to WebSocket
			outputMsg := Message{
				Type: MessageTypeTerminalOutput,
				Payload: TerminalOutput{
					PaneID: "main", // T3.x will implement proper pane management
					Data:   string(buffer[:n]),
				},
			}

			if err := conn.WriteJSON(outputMsg); err != nil {
				log.Printf("Error sending terminal output to WebSocket: %v", err)
				break
			}
		}
	}

	// Session ended - notify client
	statusMsg := Message{
		Type: MessageTypeStatus,
		Payload: StatusMessage{
			Connected: false,
			Message:   "Terminal session ended",
		},
	}
	conn.WriteJSON(statusMsg)
}

// getTmuxManager returns tmux manager for container, creating if needed
func (h *Handler) getTmuxManager(containerName string) *tmux.Manager {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if manager, exists := h.tmuxManagers[containerName]; exists {
		return manager
	}

	// Create new tmux manager for this container
	manager := tmux.NewManager(containerName)
	h.tmuxManagers[containerName] = manager
	return manager
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
