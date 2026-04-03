package main

import (
	"flag"
	"fmt"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"

	"github.com/gohive/core/logger"
	"github.com/gohive/demo-ws/config"
	"github.com/gohive/demo-ws/fxmodule"
	"github.com/gohive/demo-ws/handler"
)

func main() {
	configPath := flag.String("config", "./demo-ws/config.toml", "path to config file")
	flag.Parse()

	if err := config.Load(*configPath); err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}

	cfg := config.Cfg
	logger.Init(cfg.Log.Level, cfg.Log.Format)
	logger.Infof("Starting %s...", cfg.App.Name)

	app := fx.New(
		fx.Provide(
			fxmodule.NewHub,
			fxmodule.NewGin,
			fxmodule.NewWebSocketHandler,
		),
		fx.Invoke(registerRoutes),
	)

	app.Run()
}

func registerRoutes(r *gin.Engine, ws *handler.WebSocketHandler) {
	r.GET("/ws", ws.HandleConnection)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
}
