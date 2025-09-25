package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
	"lxc-terminal/internal/lxc"
	"lxc-terminal/internal/pty"
	"lxc-terminal/internal/tmux"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for development
	},
}

// Client represents a connected WebSocket client and its state
type Client struct {
	conn          *websocket.Conn
	containerID   string
	sessionID     string
	ptySession    *pty.Session
}

// Handler manages all WebSocket connections and routes messages
type Handler struct {
	clients    map[*websocket.Conn]*Client
	lxcManager *lxc.Manager
	mutex      sync.RWMutex
}

// NewHandler creates a new WebSocket handler
func NewHandler() *Handler {
	return &Handler{
		clients:    make(map[*websocket.Conn]*Client),
		lxcManager: lxc.NewManager(),
	}
}

// HandleWebSocket upgrades HTTP connections to WebSocket and manages the client lifecycle
func (h *Handler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}
	defer conn.Close()

	client := h.addClient(conn)
	defer h.removeClient(client)

	h.sendConnectionStatus(client, true, "Connected to Dash Terminal")

	for {
		var msg Message
		if err := conn.ReadJSON(&msg); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}
		h.handleMessage(client, &msg)
	}
}

// addClient registers a new client
func (h *Handler) addClient(conn *websocket.Conn) *Client {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	client := &Client{conn: conn}
	h.clients[conn] = client
	return client
}

// removeClient unregisters a client and cleans up resources
func (h *Handler) removeClient(client *Client) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	if client.ptySession != nil {
		client.ptySession.Close()
	}
	delete(h.clients, client.conn)
}

// handleMessage is the central router for all incoming messages
func (h *Handler) handleMessage(client *Client, msg *Message) {
	switch msg.Type {
	case MessageTypeContainerList:
		h.handleContainerList(client)
	case MessageTypeSessionList:
		h.handleSessionList(client, msg)
	case MessageTypeSessionCreate:
		h.handleSessionCreate(client, msg)
	case MessageTypeSessionSelect:
		h.handleSessionSelect(client, msg)
	case MessageTypeWindowList:
		h.handleWindowList(client, msg)
	case MessageTypeWindowCreate:
		h.handleWindowCreate(client, msg)
	case MessageTypeWindowSelect:
		h.handleWindowSelect(client, msg)
	case MessageTypePaneList:
		h.handlePaneList(client, msg)
	case MessageTypePaneCreate:
		h.handlePaneCreate(client, msg)
	case MessageTypePaneSelect:
		h.handlePaneSelect(client, msg)
	case MessageTypeTerminalInput:
		h.handleTerminalInput(client, msg)
	case MessageTypeTerminalResize:
		h.handleTerminalResize(client, msg)
	case MessageTypeContainerSelect:
		h.handleContainerSelect(client, msg)
	case MessageTypeContainerInfo:
		h.handleContainerInfo(client, msg)
	default:
		h.sendError(client, 400, "Unknown message type: "+msg.Type)
	}
}

// handleContainerList sends a list of available LXC containers to the client
func (h *Handler) handleContainerList(client *Client) {
	containers, err := h.lxcManager.ListContainers()
	if err != nil {
		h.sendError(client, 500, "Failed to list containers: "+err.Error())
		return
	}
	response := Message{Type: MessageTypeContainerList, Payload: containers}
	client.conn.WriteJSON(response)
}

// handleSessionList sends a list of tmux sessions for a given container
func (h *Handler) handleSessionList(client *Client, msg *Message) {
	var payload struct {
		ContainerName string `json:"container_name"`
	}
	if err := parsePayload(msg.Payload, &payload); err != nil {
		h.sendError(client, 400, "Invalid session list payload: "+err.Error())
		return
	}

	if payload.ContainerName == "" {
		h.sendError(client, 400, "container_name is required")
		return
	}

	// Update client state
	h.mutex.Lock()
	client.containerName = payload.ContainerName
	h.mutex.Unlock()

	tmuxManager := tmux.NewManager(client.containerName)
	sessions, err := tmuxManager.ListSessions()
	if err != nil {
		h.sendError(client, 500, "Failed to list tmux sessions: "+err.Error())
		return
	}

	// Convert to websocket.Session format
	wsSessions := make([]Session, len(sessions))
	for i, s := range sessions {
		wsSessions[i] = Session{ID: s.Name, Name: s.Name}
	}

	response := Message{Type: MessageTypeSessionList, Payload: wsSessions}
	client.conn.WriteJSON(response)
}

func (h *Handler) handleSessionCreate(client *Client, msg *Message) {
	var payload struct {
		SessionName string `json:"session_name"`
	}
	if err := parsePayload(msg.Payload, &payload); err != nil {
		h.sendError(client, 400, "Invalid session create payload: "+err.Error())
		return
	}

	if payload.SessionName == "" {
		h.sendError(client, 400, "session_name is required")
		return
	}

	h.mutex.RLock()
	containerName := client.containerName
	h.mutex.RUnlock()

	if containerName == "" {
		h.sendError(client, 400, "container_name must be set before creating a session")
		return
	}

	tmuxManager := tmux.NewManager(containerName)
	err := tmuxManager.CreateSession(payload.SessionName)
	if err != nil {
		h.sendError(client, 500, "Failed to create tmux session: "+err.Error())
		return
	}

	response := Message{Type: MessageTypeStatus, Payload: StatusMessage{Connected: true, Message: "Session created: " + payload.SessionName}}
	client.conn.WriteJSON(response)
}

func (h *Handler) handleSessionSelect(client *Client, msg *Message) {
	var payload struct {
		SessionID string `json:"session_id"`
	}
	if err := parsePayload(msg.Payload, &payload); err != nil {
		h.sendError(client, 400, "Invalid session select payload: "+err.Error())
		return
	}

	if payload.SessionID == "" {
		h.sendError(client, 400, "session_id is required")
		return
	}

	h.mutex.Lock()
	if client.containerName == "" {
		h.mutex.Unlock()
		h.sendError(client, 400, "container_name must be set before selecting a session")
		return
	}
	client.sessionID = payload.SessionID
	containerName := client.containerName
	h.mutex.Unlock()

	tmuxManager := tmux.NewManager(containerName)
	err := tmuxManager.SelectSession(payload.SessionID)
	if err != nil {
		h.sendError(client, 500, "Failed to select tmux session: "+err.Error())
		return
	}

	// For now, just send a success status.
	// In a later task, this will involve creating a PTY and forwarding I/O.
	response := Message{Type: MessageTypeStatus, Payload: StatusMessage{Connected: true, Message: "Session selected: " + payload.SessionID}}
	client.conn.WriteJSON(response)
}

func (h *Handler) handleWindowList(client *Client, msg *Message) {
	h.mutex.RLock()
	containerName := client.containerName
	sessionID := client.sessionID
	h.mutex.RUnlock()

	if containerName == "" || sessionID == "" {
		h.sendError(client, 400, "container_name and session_id must be set before listing windows")
		return
	}

	tmuxManager := tmux.NewManager(containerName)
	windows, err := tmuxManager.ListWindows(sessionID)
	if err != nil {
		h.sendError(client, 500, "Failed to list tmux windows: "+err.Error())
		return
	}

	wsWindows := make([]Window, len(windows))
	for i, w := range windows {
		wsWindows[i] = Window{ID: strconv.Itoa(w.ID), Name: w.Name, SessionID: sessionID}
	}

	response := Message{Type: MessageTypeWindowList, Payload: wsWindows}
	client.conn.WriteJSON(response)
}

func (h *Handler) handleWindowCreate(client *Client, msg *Message) {
	var payload struct {
		WindowName string `json:"window_name"`
	}
	if err := parsePayload(msg.Payload, &payload); err != nil {
		h.sendError(client, 400, "Invalid window create payload: "+err.Error())
		return
	}

	if payload.WindowName == "" {
		h.sendError(client, 400, "window_name is required")
		return
	}

	h.mutex.RLock()
	containerName := client.containerName
	sessionID := client.sessionID
	h.mutex.RUnlock()

	if containerName == "" || sessionID == "" {
		h.sendError(client, 400, "container_name and session_id must be set before creating a window")
		return
	}

	tmuxManager := tmux.NewManager(containerName)
	err := tmuxManager.CreateWindow(sessionID, payload.WindowName)
	if err != nil {
		h.sendError(client, 500, "Failed to create tmux window: "+err.Error())
		return
	}

	response := Message{Type: MessageTypeStatus, Payload: StatusMessage{Connected: true, Message: "Window created: " + payload.WindowName}}
	client.conn.WriteJSON(response)
}

func (h *Handler) handleWindowSelect(client *Client, msg *Message) {
	var payload struct {
		WindowID string `json:"window_id"`
	}
	if err := parsePayload(msg.Payload, &payload); err != nil {
		h.sendError(client, 400, "Invalid window select payload: "+err.Error())
		return
	}

	if payload.WindowID == "" {
		h.sendError(client, 400, "window_id is required")
		return
	}

	windowID, err := strconv.Atoi(payload.WindowID)
	if err != nil {
		h.sendError(client, 400, "Invalid window_id format")
		return
	}

	h.mutex.RLock()
	containerName := client.containerName
	sessionID := client.sessionID
	h.mutex.RUnlock()

	if containerName == "" || sessionID == "" {
		h.sendError(client, 400, "container_name and session_id must be set before selecting a window")
		return
	}

	tmuxManager := tmux.NewManager(containerName)
	err = tmuxManager.SelectWindow(sessionID, windowID)
	if err != nil {
		h.sendError(client, 500, "Failed to select tmux window: "+err.Error())
		return
	}

	response := Message{Type: MessageTypeStatus, Payload: StatusMessage{Connected: true, Message: "Window selected: " + payload.WindowID}}
	client.conn.WriteJSON(response)
}

func (h *Handler) handlePaneList(client *Client, msg *Message) {
	var payload struct {
		WindowID string `json:"window_id"`
	}
	if err := parsePayload(msg.Payload, &payload); err != nil {
		h.sendError(client, 400, "Invalid pane list payload: "+err.Error())
		return
	}

	if payload.WindowID == "" {
		h.sendError(client, 400, "window_id is required")
		return
	}

	windowID, err := strconv.Atoi(payload.WindowID)
	if err != nil {
		h.sendError(client, 400, "Invalid window_id format")
		return
	}

	h.mutex.RLock()
	containerName := client.containerName
	sessionID := client.sessionID
	h.mutex.RUnlock()

	if containerName == "" || sessionID == "" {
		h.sendError(client, 400, "container_name and session_id must be set before listing panes")
		return
	}

	tmuxManager := tmux.NewManager(containerName)
	panes, err := tmuxManager.ListPanes(sessionID, windowID)
	if err != nil {
		h.sendError(client, 500, "Failed to list tmux panes: "+err.Error())
		return
	}

	wsPanes := make([]Pane, len(panes))
	for i, p := range panes {
		wsPanes[i] = Pane{ID: strconv.Itoa(p.ID), WindowID: payload.WindowID, Width: p.Width, Height: p.Height}
	}

	response := Message{Type: MessageTypePaneList, Payload: wsPanes}
	client.conn.WriteJSON(response)
}

func (h *Handler) handlePaneCreate(client *Client, msg *Message) {
	var payload struct {
		WindowID string `json:"window_id"`
	}
	if err := parsePayload(msg.Payload, &payload); err != nil {
		h.sendError(client, 400, "Invalid pane create payload: "+err.Error())
		return
	}

	if payload.WindowID == "" {
		h.sendError(client, 400, "window_id is required")
		return
	}

	windowID, err := strconv.Atoi(payload.WindowID)
	if err != nil {
		h.sendError(client, 400, "Invalid window_id format")
		return
	}

	h.mutex.RLock()
	containerName := client.containerName
	sessionID := client.sessionID
	h.mutex.RUnlock()

	if containerName == "" || sessionID == "" {
		h.sendError(client, 400, "container_name and session_id must be set before creating a pane")
		return
	}

	tmuxManager := tmux.NewManager(containerName)
	err = tmuxManager.CreatePane(sessionID, windowID)
	if err != nil {
		h.sendError(client, 500, "Failed to create tmux pane: "+err.Error())
		return
	}

	response := Message{Type: MessageTypeStatus, Payload: StatusMessage{Connected: true, Message: "Pane created in window: " + payload.WindowID}}
	client.conn.WriteJSON(response)
}

func (h *Handler) handlePaneSelect(client *Client, msg *Message) {
	var payload struct {
		PaneID string `json:"pane_id"`
	}
	if err := parsePayload(msg.Payload, &payload); err != nil {
		h.sendError(client, 400, "Invalid pane select payload: "+err.Error())
		return
	}

	if payload.PaneID == "" {
		h.sendError(client, 400, "pane_id is required")
		return
	}

	parts := strings.Split(payload.PaneID, ".")
	if len(parts) != 2 {
		h.sendError(client, 400, "Invalid pane_id format. Expected windowID.paneID")
		return
	}
	windowID, err := strconv.Atoi(parts[0])
	if err != nil {
		h.sendError(client, 400, "Invalid windowID in pane_id")
		return
	}
	paneID, err := strconv.Atoi(parts[1])
	if err != nil {
		h.sendError(client, 400, "Invalid paneID in pane_id")
		return
	}

	h.mutex.RLock()
	containerName := client.containerName
	sessionID := client.sessionID
	h.mutex.RUnlock()

	if containerName == "" || sessionID == "" {
		h.sendError(client, 400, "container_name and session_id must be set before selecting a pane")
		return
	}

	tmuxManager := tmux.NewManager(containerName)
	err = tmuxManager.SelectPane(sessionID, windowID, paneID)
	if err != nil {
		h.sendError(client, 500, "Failed to select tmux pane: "+err.Error())
		return
	}

	response := Message{Type: MessageTypeStatus, Payload: StatusMessage{Connected: true, Message: "Pane selected: " + payload.PaneID}}
	client.conn.WriteJSON(response)
}

func (h *Handler) handleTerminalInput(client *Client, msg *Message) {
	var payload TerminalInput
	if err := parsePayload(msg.Payload, &payload); err != nil {
		h.sendError(client, 400, "Invalid terminal input payload: "+err.Error())
		return
	}

	h.mutex.RLock()
	containerName := client.containerName
	sessionID := client.sessionID
	h.mutex.RUnlock()

	if containerName == "" || sessionID == "" {
		h.sendError(client, 400, "container_name and session_id must be set before sending terminal input")
		return
	}

	// Assuming paneID is in the format "windowID.paneID"
	parts := strings.Split(payload.PaneID, ".")
	if len(parts) != 2 {
		h.sendError(client, 400, "Invalid pane_id format. Expected windowID.paneID")
		return
	}
	windowID, err := strconv.Atoi(parts[0])
	if err != nil {
		h.sendError(client, 400, "Invalid windowID in pane_id")
		return
	}
	paneID, err := strconv.Atoi(parts[1])
	if err != nil {
		h.sendError(client, 400, "Invalid paneID in pane_id")
		return
	}

	tmuxManager := tmux.NewManager(containerName)
	err = tmuxManager.SendKeys(sessionID, windowID, paneID, payload.Data)
	if err != nil {
		h.sendError(client, 500, "Failed to send keys to tmux: "+err.Error())
		return
	}
}

func (h *Handler) handleTerminalResize(client *Client, msg *Message) {
	var payload TerminalResize
	if err := parsePayload(msg.Payload, &payload); err != nil {
		h.sendError(client, 400, "Invalid terminal resize payload: "+err.Error())
		return
	}

	h.mutex.RLock()
	containerName := client.containerName
	sessionID := client.sessionID
	h.mutex.RUnlock()

	if containerName == "" || sessionID == "" {
		h.sendError(client, 400, "container_name and session_id must be set before resizing terminal")
		return
	}

	// Assuming paneID is in the format "windowID.paneID"
	parts := strings.Split(payload.PaneID, ".")
	if len(parts) != 2 {
		h.sendError(client, 400, "Invalid pane_id format. Expected windowID.paneID")
		return
	}
	windowID, err := strconv.Atoi(parts[0])
	if err != nil {
		h.sendError(client, 400, "Invalid windowID in pane_id")
		return
	}
	paneID, err := strconv.Atoi(parts[1])
	if err != nil {
		h.sendError(client, 400, "Invalid paneID in pane_id")
		return
	}

	tmuxManager := tmux.NewManager(containerName)
	err = tmuxManager.ResizePane(sessionID, windowID, paneID, payload.Width, payload.Height)
	if err != nil {
		h.sendError(client, 500, "Failed to resize tmux pane: "+err.Error())
		return
	}
}

func (h *Handler) handleContainerSelect(client *Client, msg *Message) {
	// Container selection is handled by sending a session_list message
	// with the desired container_name. This sets the container context for the client.
	h.sendError(client, 501, "Message type '"+msg.Type+"' is deprecated. Use 'session_list'.")
}

func (h *Handler) handleContainerInfo(client *Client, msg *Message) {
	var payload ContainerInfoMessage
	if err := parsePayload(msg.Payload, &payload); err != nil {
		h.sendError(client, 400, "Invalid container info payload: "+err.Error())
		return
	}

	if payload.ContainerName == "" {
		h.sendError(client, 400, "container_name is required")
		return
	}

	container, err := h.lxcManager.GetContainer(payload.ContainerName)
	if err != nil {
		h.sendError(client, 500, "Failed to get container info: "+err.Error())
		return
	}

	response := Message{Type: MessageTypeContainerInfo, Payload: container}
	client.conn.WriteJSON(response)
}


// --- Utility Functions ---

func (h *Handler) sendConnectionStatus(client *Client, connected bool, message string) {
	statusMsg := Message{
		Type:    MessageTypeStatus,
		Payload: StatusMessage{Connected: connected, Message: message},
	}
	client.conn.WriteJSON(statusMsg)
}

func (h *Handler) sendError(client *Client, code int, message string) {
	errorMsg := Message{
		Type:    MessageTypeError,
		Payload: ErrorMessage{Code: code, Message: message},
	}
	client.conn.WriteJSON(errorMsg)
}

func parsePayload(payload interface{}, target interface{}) error {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return json.Unmarshal(payloadBytes, target)
}

// Message represents a WebSocket message structure
type Message struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}