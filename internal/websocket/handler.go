package websocket

import (
	"encoding/json"
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
		return true // Allow all origins for development
	},
}

// Client represents a connected WebSocket client and its state
type Client struct {
	conn          *websocket.Conn
	containerName string
	sessionName   string
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
	// Add other message handlers here based on t2-1_task_description.yaml
	// For now, we will add placeholders
	case MessageTypeSessionCreate, MessageTypeSessionSelect, MessageTypeWindowList, MessageTypeWindowCreate, MessageTypeWindowSelect, MessageTypePaneList, MessageTypePaneCreate, MessageTypePaneSelect, MessageTypeTerminalInput, MessageTypeTerminalResize, MessageTypeContainerSelect, MessageTypeContainerInfo:
		h.sendError(client, 501, "Message type '"+msg.Type+"' not yet implemented.")
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