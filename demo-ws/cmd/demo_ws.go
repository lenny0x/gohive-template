package cmd

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/gohive/core/logger"
	"github.com/gohive/demo-ws/config"
	"github.com/gohive/demo-ws/handler"
	"github.com/gohive/demo-ws/hub"
)

func Run(configPath string) {
	// Load config
	if err := config.Load(configPath); err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}
	cfg := config.Cfg

	// Init logger
	logger.Init(cfg.Log.Level, cfg.Log.Format)
	logger.Infof("Starting %s...", cfg.App.Name)

	// Create hub
	h := hub.NewHub()
	go h.Run()

	// Setup router
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	wsHandler := handler.NewWebSocketHandler(
		h,
		cfg.WebSocket.ReadBufferSize,
		cfg.WebSocket.WriteBufferSize,
		cfg.WebSocket.AllowedOrigins,
	)
	r.GET("/ws", wsHandler.HandleConnection)

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Start server
	addr := fmt.Sprintf(":%d", cfg.App.Port)
	logger.Infof("Server listening on %s", addr)
	if err := r.Run(addr); err != nil {
		logger.Fatalf("failed to start server: %v", err)
	}
}
