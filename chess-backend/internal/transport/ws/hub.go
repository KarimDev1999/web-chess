package ws

import (
	"encoding/json"
	"sync"

	"chess-backend/internal/transport/wsmsg"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID     string
	UserID string
	Conn   *websocket.Conn
	Send   chan []byte
	Hub    *Hub
}

type Hub struct {
	clients     map[string]*Client
	userClients map[string][]string
	gameRooms   map[string]map[string]bool
	broadcast   chan []byte
	unregister  chan *Client
	mu          sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients:     make(map[string]*Client),
		userClients: make(map[string][]string),
		gameRooms:   make(map[string]map[string]bool),
		broadcast:   make(chan []byte),
		unregister:  make(chan *Client),
	}
}

func (h *Hub) Register(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[client.ID] = client
	h.userClients[client.UserID] = append(h.userClients[client.UserID], client.ID)
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.unregister:
			h.mu.Lock()
			delete(h.clients, client.ID)

			conns := h.userClients[client.UserID]
			for i, id := range conns {
				if id == client.ID {
					h.userClients[client.UserID] = append(conns[:i], conns[i+1:]...)
					break
				}
			}
			if len(h.userClients[client.UserID]) == 0 {
				delete(h.userClients, client.UserID)
			}

			h.removeFromAllGames(client.UserID)
			close(client.Send)
			h.mu.Unlock()

		case message := <-h.broadcast:
			h.mu.RLock()
			for _, client := range h.clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (h *Hub) SendToUser(userID string, message []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	conns := h.userClients[userID]
	for _, connID := range conns {
		client, ok := h.clients[connID]
		if ok {
			select {
			case client.Send <- message:
			default:

			}
		}
	}
}

func (h *Hub) NotifyPlayers(playerIDs []string, data []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for _, pid := range playerIDs {
		if pid == "" {
			continue
		}
		conns := h.userClients[pid]
		for _, connID := range conns {
			client, ok := h.clients[connID]
			if ok {
				select {
				case client.Send <- data:
				default:

				}
			}
		}
	}
}

func (h *Hub) JoinGame(userID, gameID string) {
	h.mu.Lock()
	if h.gameRooms[gameID] == nil {
		h.gameRooms[gameID] = make(map[string]bool)
	}
	h.gameRooms[gameID][userID] = true
	h.mu.Unlock()

	h.broadcastPresence(gameID)
}

func (h *Hub) LeaveGame(userID, gameID string) {
	h.mu.Lock()
	if room := h.gameRooms[gameID]; room != nil {
		delete(room, userID)
		if len(room) == 0 {
			delete(h.gameRooms, gameID)
		}
	}
	h.mu.Unlock()

	h.broadcastPresence(gameID)
}

func (h *Hub) removeFromAllGames(userID string) {
	for gameID, room := range h.gameRooms {
		delete(room, userID)
		if len(room) == 0 {
			delete(h.gameRooms, gameID)
		}
	}
}

func (h *Hub) GetPresence(gameID string) map[string]bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	room := h.gameRooms[gameID]
	if room == nil {
		return nil
	}

	result := make(map[string]bool)
	for uid := range room {
		result[uid] = true
	}
	return result
}

func (h *Hub) broadcastPresence(gameID string) {
	room := h.gameRooms[gameID]
	if room == nil {
		return
	}

	userIDs := make([]string, 0, len(room))
	for uid := range room {
		userIDs = append(userIDs, uid)
	}

	msg, err := json.Marshal(wsmsg.NewPresence(gameID, userIDs))
	if err != nil {
		return
	}

	for uid := range room {
		h.SendToUser(uid, msg)
	}
}
