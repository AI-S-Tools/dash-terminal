package websocket

import (
	"sync"

	"github.com/gorilla/websocket"
)

// Client represents a single WebSocket client connection
type Client struct {
	Conn          *websocket.Conn
	ContainerName string
}

// ClientManager manages all active WebSocket clients
type ClientManager struct {
	clients map[*websocket.Conn]*Client
	mutex   sync.RWMutex
}

// NewClientManager creates a new ClientManager
func NewClientManager() *ClientManager {
	return &ClientManager{
		clients: make(map[*websocket.Conn]*Client),
	}
}

// AddClient adds a new client to the manager
func (m *ClientManager) AddClient(conn *websocket.Conn) *Client {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	client := &Client{Conn: conn}
	m.clients[conn] = client
	return client
}

// RemoveClient removes a client from the manager
func (m.ClientManager) RemoveClient(conn *websocket.Conn) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	delete(m.clients, conn)
}

// GetClient retrieves a client from the manager
func (m *ClientManager) GetClient(conn *websocket.Conn) (*Client, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	client, ok := m.clients[conn]
	return client, ok
}
