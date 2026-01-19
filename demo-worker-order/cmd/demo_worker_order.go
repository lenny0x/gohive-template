package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gohive/core/logger"
	"github.com/gohive/demo-worker-order/config"
	"github.com/gohive/demo-worker-order/task"
	database "github.com/gohive/pkg/mysql"
	"github.com/gohive/pkg/redis"
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

	// Init redis
	if err := redis.Init(redis.Config{
		Host:     cfg.Redis.Host,
		Port:     cfg.Redis.Port,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	}); err != nil {
		logger.Warnf("failed to init redis: %v", err)
	}

	// Register tasks
	runner := task.NewTaskRunner()
	runner.Register(task.NewOrderExpireTask())
	runner.Register(task.NewCleanupTask())

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start worker pool
	var wg sync.WaitGroup
	taskChan := make(chan task.Task, cfg.Worker.Concurrency)

	// Start workers
	for i := 0; i < cfg.Worker.Concurrency; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for t := range taskChan {
				logger.Debugf("worker %d executing task: %s", workerID, t.Name())
				if err := t.Run(ctx); err != nil {
					logger.Errorf("task %s failed: %v", t.Name(), err)
				}
			}
		}(i)
	}

	// Ticker to schedule tasks
	ticker := time.NewTicker(time.Duration(cfg.Worker.Interval) * time.Second)
	defer ticker.Stop()

	logger.Infof("Worker started with %d concurrent workers, interval: %ds",
		cfg.Worker.Concurrency, cfg.Worker.Interval)

	// Signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Main loop
	for {
		select {
		case <-ticker.C:
			for _, t := range runner.Tasks() {
				select {
				case taskChan <- t:
				default:
					logger.Warnf("task channel full, skipping task: %s", t.Name())
				}
			}
		case <-sigChan:
			logger.Info("Shutting down...")
			cancel()
			close(taskChan)
			wg.Wait()
			return
		}
	}
}
