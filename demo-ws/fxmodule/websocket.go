package fxmodule

import (
	"github.com/gohive/demo-ws/config"
	"github.com/gohive/demo-ws/handler"
	"github.com/gohive/demo-ws/hub"
)

func NewWebSocketHandler(h *hub.Hub) *handler.WebSocketHandler {
	cfg := config.Cfg.WebSocket
	return handler.NewWebSocketHandler(
		h,
		cfg.ReadBufferSize,
		cfg.WriteBufferSize,
		cfg.AllowedOrigins,
	)
}
