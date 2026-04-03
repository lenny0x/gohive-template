package fxmodule

import (
	"context"
	"sync"
	"time"

	"go.uber.org/fx"

	"github.com/gohive/core/logger"
	"github.com/gohive/demo-worker-order/config"
	"github.com/gohive/demo-worker-order/task"
)

func NewTaskRunner() *task.TaskRunner {
	runner := task.NewTaskRunner()
	runner.Register(task.NewOrderExpireTask())
	runner.Register(task.NewCleanupTask())
	return runner
}

func StartWorker(lc fx.Lifecycle, runner *task.TaskRunner) {
	cfg := config.Cfg.Worker
	ctx, cancel := context.WithCancel(context.Background())
	taskChan := make(chan task.Task, cfg.Concurrency)
	var wg sync.WaitGroup

	lc.Append(fx.Hook{
		OnStart: func(startCtx context.Context) error {
			for i := 0; i < cfg.Concurrency; i++ {
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

			ticker := time.NewTicker(time.Duration(cfg.Interval) * time.Second)
			go func() {
				defer ticker.Stop()
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
					case <-ctx.Done():
						return
					}
				}
			}()

			logger.Infof("Worker started with %d concurrent workers, interval: %ds",
				cfg.Concurrency, cfg.Interval)
			return nil
		},
		OnStop: func(stopCtx context.Context) error {
			logger.Info("Shutting down worker...")
			cancel()
			close(taskChan)
			wg.Wait()
			return nil
		},
	})
}
