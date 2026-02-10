package handlers

import (
	"net/http"

	"github.com/Mahaveer86619/Hearth/pkg/logger"
	"github.com/Mahaveer86619/Hearth/pkg/services"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for development
	},
}

type WSHandler struct {
}

func NewWSHandler() *WSHandler {
	return &WSHandler{}
}

func (h *WSHandler) Handle(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Error("WS", "Failed to upgrade to websocket: %v", err)
		return
	}

	logger.Info("WS", "New WebSocket connection from %s", conn.RemoteAddr().String())

	hub := services.GetHub()
	if hub == nil {
		logger.Error("WS", "Hub not initialized")
		conn.Close()
		return
	}

	client := &services.Client{Hub: hub, Send: make(chan []byte, 256)}
	client.Hub.Register <- client

	go client.WritePump(conn)
	go client.ReadPump(conn)
}
