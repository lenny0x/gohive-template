package fxmodule

import (
	"context"

	"go.uber.org/fx"

	"github.com/gohive/core/logger"
	"github.com/gohive/demo-ws/hub"
)

func NewHub(lc fx.Lifecycle) *hub.Hub {
	h := hub.NewHub()

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go h.Run()
			logger.Info("WebSocket hub started")
			return nil
		},
	})

	return h
}
