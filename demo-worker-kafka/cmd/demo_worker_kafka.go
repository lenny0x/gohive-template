package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gohive/core/logger"
	"github.com/gohive/demo-worker-kafka/config"
	"github.com/gohive/demo-worker-kafka/consumer"
	"github.com/gohive/demo-worker-kafka/handler"
	"github.com/gohive/pkg/kafka"
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

	kafkaCfg := kafka.Config{
		Brokers: cfg.Kafka.Brokers,
		GroupID: cfg.Kafka.GroupID,
	}

	// Create consumer and register handlers
	c := consumer.NewConsumer()
	c.RegisterHandler("orders", kafka.NewReader(kafkaCfg, "orders").Reader, handler.NewOrderHandler())
	c.RegisterHandler("notifications", kafka.NewReader(kafkaCfg, "notifications").Reader, handler.NewNotificationHandler())

	// Start consuming
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := c.Start(ctx); err != nil {
			logger.Errorf("consumer error: %v", err)
		}
	}()

	logger.Infof("Kafka consumer started, listening on topics: %v", cfg.Topics)

	// Wait for shutdown signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	logger.Info("Shutting down...")
	cancel()
	c.Close()
}
