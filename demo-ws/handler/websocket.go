package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gohive/core/logger"
	"github.com/gohive/demo-ws/hub"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type WebSocketHandler struct {
	hub            *hub.Hub
	upgrader       websocket.Upgrader
	allowedOrigins []string
}

func NewWebSocketHandler(h *hub.Hub, readBufferSize, writeBufferSize int, allowedOrigins []string) *WebSocketHandler {
	handler := &WebSocketHandler{
		hub:            h,
		allowedOrigins: allowedOrigins,
	}
	handler.upgrader = websocket.Upgrader{
		ReadBufferSize:  readBufferSize,
		WriteBufferSize: writeBufferSize,
		CheckOrigin:     handler.checkOrigin,
	}
	return handler
}

func (h *WebSocketHandler) checkOrigin(r *http.Request) bool {
	if len(h.allowedOrigins) == 0 {
		return true
	}
	origin := r.Header.Get("Origin")
	for _, allowed := range h.allowedOrigins {
		if allowed == "*" || allowed == origin {
			return true
		}
	}
	return false
}

func (h *WebSocketHandler) HandleConnection(c *gin.Context) {
	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Errorf("failed to upgrade connection: %v", err)
		return
	}

	clientID := uuid.New().String()
	client := &hub.Client{
		ID:   clientID,
		Conn: conn,
		Send: make(chan []byte, 256),
	}

	h.hub.Register(client)
	logger.Infof("client connected: %s", clientID)

	go h.writePump(client)
	go h.readPump(client)
}

func (h *WebSocketHandler) readPump(client *hub.Client) {
	defer func() {
		h.hub.Unregister(client)
		client.Conn.Close()
		logger.Infof("client disconnected: %s", client.ID)
	}()

	for {
		_, message, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logger.Errorf("read error: %v", err)
			}
			break
		}

		// Handle message
		logger.Debugf("received message from %s: %s", client.ID, string(message))
		h.hub.Broadcast(message)
	}
}

func (h *WebSocketHandler) writePump(client *hub.Client) {
	defer client.Conn.Close()

	for message := range client.Send {
		if err := client.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
			logger.Errorf("write error: %v", err)
			return
		}
	}
}
