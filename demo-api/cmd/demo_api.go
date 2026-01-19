package cmd

import (
	"fmt"

	"github.com/gohive/core/logger"
	"github.com/gohive/demo-api/config"
	"github.com/gohive/demo-api/router"
	database "github.com/gohive/pkg/mysql"
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

	// Init database
	if err := database.Init(database.Config{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.DBName,
	}); err != nil {
		logger.Fatalf("failed to init database: %v", err)
	}

	// Setup router and start server
	r := router.Setup()
	addr := fmt.Sprintf(":%d", cfg.App.Port)
	logger.Infof("Server listening on %s", addr)
	if err := r.Run(addr); err != nil {
		logger.Fatalf("failed to start server: %v", err)
	}
}
