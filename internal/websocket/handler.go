package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"lxc-terminal/internal/lxc"
	"lxc-terminal/internal/pty"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for development
	},
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
	client := &Client{Conn: conn}
	h.clients[conn] = client
	return client
}

// removeClient unregisters a client and cleans up resources
func (h *Handler) removeClient(client *Client) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	if client.PtySession != nil {
		client.PtySession.Close()
	}
	delete(h.clients, client.Conn)
}

// handleMessage is the central router for all incoming messages
func (h *Handler) handleMessage(client *Client, msg *Message) {
	switch msg.Type {
	case MessageTypeContainerList:
		h.handleContainerList(client)
	case MessageTypeContainerSelect:
		h.handleContainerSelect(client, msg)
	case MessageTypeTerminalInput:
		h.handleTerminalInput(client, msg)
	case MessageTypeTerminalResize:
		h.handleTerminalResize(client, msg)
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
	client.Conn.WriteJSON(response)
}

// handleContainerSelect starts a PTY session in the selected container
func (h *Handler) handleContainerSelect(client *Client, msg *Message) {
	var payload ContainerSelectMessage
	if err := parsePayload(msg.Payload, &payload); err != nil {
		h.sendError(client, 400, "Invalid container select payload: "+err.Error())
		return
	}

	if payload.ContainerName == "" {
		h.sendError(client, 400, "container_name is required")
		return
	}

	h.mutex.Lock()
	client.ContainerName = payload.ContainerName
	h.mutex.Unlock()

	ptySession, err := pty.NewSession(payload.ContainerName)
	if err != nil {
		h.sendError(client, 500, "Failed to start PTY session: "+err.Error())
		return
	}

	h.mutex.Lock()
	client.PtySession = ptySession
	h.mutex.Unlock()

	go h.streamPtyOutput(client, ptySession)

	h.sendConnectionStatus(client, true, "PTY session started in "+payload.ContainerName)
}

// streamPtyOutput forwards output from the PTY to the WebSocket client
func (h *Handler) streamPtyOutput(client *Client, ptySession *pty.Session) {
	buf := make([]byte, 1024)
	for {
		n, err := ptySession.Read(buf)
		if err != nil {
			// Session closed
			break
		}
		output := string(buf[:n])
		response := Message{Type: MessageTypeTerminalOutput, Payload: TerminalOutput{Data: output}}
		if err := client.Conn.WriteJSON(response); err != nil {
			// Client disconnected
			break
		}
	}
}

func (h *Handler) handleTerminalInput(client *Client, msg *Message) {
	var payload TerminalInput
	if err := parsePayload(msg.Payload, &payload); err != nil {
		h.sendError(client, 400, "Invalid terminal input payload: "+err.Error())
		return
	}

	h.mutex.RLock()
	ptySession := client.PtySession
	h.mutex.RUnlock()

	if ptySession == nil {
		h.sendError(client, 400, "No active PTY session")
		return
	}

	if _, err := ptySession.Write([]byte(payload.Data)); err != nil {
		h.sendError(client, 500, "Failed to write to PTY: "+err.Error())
	}
}

func (h *Handler) handleTerminalResize(client *Client, msg *Message) {
	var payload TerminalResize
	if err := parsePayload(msg.Payload, &payload); err != nil {
		h.sendError(client, 400, "Invalid terminal resize payload: "+err.Error())
		return
	}

	h.mutex.RLock()
	ptySession := client.PtySession
	h.mutex.RUnlock()

	if ptySession == nil {
		h.sendError(client, 400, "No active PTY session")
		return
	}

	if err := ptySession.Resize(payload.Width, payload.Height); err != nil {
		h.sendError(client, 500, "Failed to resize PTY: "+err.Error())
	}
}

// --- Utility Functions ---

func (h *Handler) sendConnectionStatus(client *Client, connected bool, message string) {
	statusMsg := Message{
		Type:    MessageTypeStatus,
		Payload: StatusMessage{Connected: connected, Message: message},
	}
	client.Conn.WriteJSON(statusMsg)
}

func (h *Handler) sendError(client *Client, code int, message string) {
	errorMsg := Message{
		Type:    MessageTypeError,
		Payload: ErrorMessage{Code: code, Message: message},
	}
	client.Conn.WriteJSON(errorMsg)
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