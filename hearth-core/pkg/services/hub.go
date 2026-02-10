package services

import (
	"context"
	"sync"
	"time"

	"github.com/Mahaveer86619/Hearth/pkg/constants"
	"github.com/Mahaveer86619/Hearth/pkg/db"
	"github.com/Mahaveer86619/Hearth/pkg/logger"
	"github.com/gorilla/websocket"
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

// Hub maintains the set of active clients and broadcasts messages to them.
type Hub struct {
	// Registered clients map.
	clients map[*Client]bool

	// Inbound messages from Redis.
	broadcast chan []byte

	// Register requests from the clients.
	Register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	mu sync.RWMutex
}

// Client acts as an intermediary between the websocket connection and the hub.
type Client struct {
	Hub  *Hub
	Send chan []byte
}

var globalHub *Hub

func NewHub() *Hub {
	globalHub = &Hub{
		broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
	return globalHub
}

func GetHub() *Hub {
	return globalHub
}

func (h *Hub) Run() {
	logger.Info("Hub", "Broadcast Hub started")

	go h.subscribeToRedis()

	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			logger.Info("Hub", "Client registered. Total: %d", len(h.clients))

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.Send)
			}
			h.mu.Unlock()
			logger.Info("Hub", "Client unregistered. Total: %d", len(h.clients))

		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.Send <- message:
				default:
					// If client buffer is full, close and drop to prevent blocking the Hub
					close(client.Send)
					delete(h.clients, client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (h *Hub) subscribeToRedis() {
	if db.RDB == nil {
		logger.Error("Hub", "Redis client not initialized")
		return
	}

	pubsub := db.RDB.Subscribe(context.Background(), string(constants.RedisChannelLiveLogs))
	defer pubsub.Close()

	ch := pubsub.Channel()

	for msg := range ch {
		h.broadcast <- []byte(msg.Payload)
	}
}

func (c *Client) ReadPump(conn *websocket.Conn) {
	defer func() {
		c.Hub.unregister <- c
		conn.Close()
	}()
	conn.SetReadLimit(512)
	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error { conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
		// We ignore incoming messages for now (Pulse is one-way: Server -> Client)
	}
}

func (c *Client) WritePump(conn *websocket.Conn) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
