package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/websocket"
	"lxc-terminal/internal/zellij"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for development
	},
}

// Handler manages API endpoints for Zellij daemon
type Handler struct {
	daemon *zellij.Daemon
}

// NewHandler creates a new API handler
func NewHandler(daemon *zellij.Daemon) *Handler {
	return &Handler{daemon: daemon}
}

// HandleWebSocket upgrades to WebSocket and bridges PTY
func (h *Handler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session")
	if sessionID == "" {
		http.Error(w, "session parameter required", http.StatusBadRequest)
		return
	}

	// Upgrade to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}
	defer conn.Close()

	// Attach to Zellij session
	ptyFile, err := h.daemon.AttachToSession(sessionID)
	if err != nil {
		log.Printf("Failed to attach to session %s: %v", sessionID, err)
		conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Error: %v", err)))
		return
	}
	defer ptyFile.Close()

	// Bridge PTY ↔ WebSocket
	go h.bridgePTYToWS(ptyFile, conn)
	go h.bridgeWSToPTY(conn, ptyFile)

	// Keep connection alive - input/output is handled by bridge functions
	select {}
}

// bridgePTYToWS reads from PTY and sends to WebSocket
func (h *Handler) bridgePTYToWS(ptyFile *os.File, conn *websocket.Conn) {
	buffer := make([]byte, 1024)
	for {
		n, err := ptyFile.Read(buffer)
		if err != nil {
			if err != io.EOF {
				log.Printf("PTY read error: %v", err)
			}
			break
		}

		if err := conn.WriteMessage(websocket.BinaryMessage, buffer[:n]); err != nil {
			log.Printf("WebSocket write error: %v", err)
			break
		}
	}
}

// bridgeWSToPTY reads from WebSocket and sends to PTY
func (h *Handler) bridgeWSToPTY(conn *websocket.Conn, ptyFile *os.File) {
	for {
		msgType, data, err := conn.ReadMessage()
		if err != nil {
			log.Printf("WebSocket read error: %v", err)
			break
		}

		if msgType == websocket.TextMessage {
			// Check if it's a control message (JSON)
			if len(data) > 0 && data[0] == '{' {
				var msg map[string]interface{}
				if err := json.Unmarshal(data, &msg); err == nil {
					// Handle control messages like resize
					if msgType, ok := msg["type"].(string); ok && msgType == "resize" {
						if width, ok := msg["width"].(float64); ok {
							if height, ok := msg["height"].(float64); ok {
								if err := zellij.ResizePTY(ptyFile, int(width), int(height)); err != nil {
									log.Printf("Resize error: %v", err)
								}
							}
						}
					}
					continue // Don't send control messages to PTY
				}
			}

			// Regular terminal input - send to PTY
			if _, err := ptyFile.Write(data); err != nil {
				log.Printf("PTY write error: %v", err)
				break
			}
		} else if msgType == websocket.BinaryMessage {
			// Binary terminal data
			if _, err := ptyFile.Write(data); err != nil {
				log.Printf("PTY write error: %v", err)
				break
			}
		}
	}
}

// HandleSessions handles session list and creation
func (h *Handler) HandleSessions(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.listSessions(w, r)
	case http.MethodPost:
		h.createSession(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// HandleSession handles individual session operations
func (h *Handler) HandleSession(w http.ResponseWriter, r *http.Request) {
	// Extract session ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/sessions/")
	parts := strings.Split(path, "/")
	if len(parts) < 1 {
		http.Error(w, "Session ID required", http.StatusBadRequest)
		return
	}

	sessionID := parts[0]

	switch r.Method {
	case http.MethodDelete:
		h.deleteSession(w, r, sessionID)
	case http.MethodGet:
		if len(parts) > 1 && parts[1] == "info" {
			h.getSessionInfo(w, r, sessionID)
		} else {
			h.getSession(w, r, sessionID)
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// listSessions returns all sessions
func (h *Handler) listSessions(w http.ResponseWriter, r *http.Request) {
	sessions := h.daemon.ListSessions()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sessions)
}

// createSession creates a new session
func (h *Handler) createSession(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		req.Name = "default"
	}

	session, err := h.daemon.CreateSession(req.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(session)
}

// getSession returns session details
func (h *Handler) getSession(w http.ResponseWriter, r *http.Request, sessionID string) {
	session, exists := h.daemon.GetSession(sessionID)
	if !exists {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(session)
}

// deleteSession removes a session
func (h *Handler) deleteSession(w http.ResponseWriter, r *http.Request, sessionID string) {
	if err := h.daemon.DeleteSession(sessionID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// getSessionInfo returns detailed session information
func (h *Handler) getSessionInfo(w http.ResponseWriter, r *http.Request, sessionID string) {
	info, err := h.daemon.GetSessionInfo(sessionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(info)
}

// HandleHealth returns daemon health status
func (h *Handler) HandleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	health := map[string]interface{}{
		"status": "ok",
		"container": h.daemon.GetContainer(),
		"sessions": len(h.daemon.ListSessions()),
	}

	// Check if Zellij is available
	if version, err := h.daemon.GetZellijVersion(); err == nil {
		health["zellij_version"] = version
		health["zellij_available"] = true
	} else {
		health["zellij_available"] = false
		health["zellij_error"] = err.Error()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(health)
}