package fxmodule

import (
	"context"

	"go.uber.org/fx"

	"github.com/gohive/core/logger"
	"github.com/gohive/demo-worker-kafka/config"
	"github.com/gohive/demo-worker-kafka/consumer"
	"github.com/gohive/demo-worker-kafka/handler"
	"github.com/gohive/pkg/kafka"
)

func NewConsumer(lc fx.Lifecycle) {
	cfg := config.Cfg
	kafkaCfg := kafka.Config{
		Brokers: cfg.Kafka.Brokers,
		GroupID: cfg.Kafka.GroupID,
	}

	c := consumer.NewConsumer()
	c.RegisterHandler("orders", kafka.NewReader(kafkaCfg, "orders").Reader, handler.NewOrderHandler())
	c.RegisterHandler("notifications", kafka.NewReader(kafkaCfg, "notifications").Reader, handler.NewNotificationHandler())

	ctx, cancel := context.WithCancel(context.Background())

	lc.Append(fx.Hook{
		OnStart: func(startCtx context.Context) error {
			go func() {
				if err := c.Start(ctx); err != nil && ctx.Err() == nil {
					logger.Errorf("consumer error: %v", err)
				}
			}()
			logger.Infof("Kafka consumer started, listening on topics: %v", cfg.Topics)
			return nil
		},
		OnStop: func(stopCtx context.Context) error {
			logger.Info("Shutting down kafka consumer...")
			cancel()
			return c.Close()
		},
	})

}
