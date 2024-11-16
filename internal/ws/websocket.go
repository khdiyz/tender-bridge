package ws

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Notification struct {
	UserID  string `json:"user_id"`
	Message string `json:"message"`
}

type WebSocketHub struct {
	clients    map[string]*websocket.Conn // userID => connection
	broadcast  chan Notification
	register   chan subscription
	unregister chan subscription
	mu         sync.Mutex
}

type subscription struct {
	conn   *websocket.Conn
	userID string
}

var hub = WebSocketHub{
	clients:    make(map[string]*websocket.Conn),
	broadcast:  make(chan Notification),
	register:   make(chan subscription),
	unregister: make(chan subscription),
}

func (h *WebSocketHub) Run() {
	for {
		select {
		case sub := <-h.register:
			h.mu.Lock()
			h.clients[sub.userID] = sub.conn
			h.mu.Unlock()
		case sub := <-h.unregister:
			h.mu.Lock()
			if conn, ok := h.clients[sub.userID]; ok && conn == sub.conn {
				delete(h.clients, sub.userID)
				sub.conn.Close()
			}
			h.mu.Unlock()
		case notification := <-h.broadcast:
			h.mu.Lock()
			if conn, ok := h.clients[notification.UserID]; ok {
				err := conn.WriteJSON(notification)
				if err != nil {
					conn.Close()
					delete(h.clients, notification.UserID)
				}
			}
			h.mu.Unlock()
		}
	}
}

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "Missing user_id", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %s", err)
		return
	}

	sub := subscription{conn: conn, userID: userID}
	hub.register <- sub

	defer func() {
		hub.unregister <- sub
	}()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

func BroadcastNotification(userID, message string) {
	hub.broadcast <- Notification{
		UserID:  userID,
		Message: message,
	}
}

func StartWebSocketHub() {
	go hub.Run()
}
