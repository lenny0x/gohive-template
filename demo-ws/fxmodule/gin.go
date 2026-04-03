package fxmodule

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"

	"github.com/gohive/core/config"
	"github.com/gohive/core/logger"
	"github.com/gohive/core/middleware"
	localConfig "github.com/gohive/demo-ws/config"
)

func NewGin(lc fx.Lifecycle) *gin.Engine {
	cfg := localConfig.Cfg.App

	if config.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(middleware.RequestID())
	r.Use(middleware.Recovery())
	r.Use(gin.Logger())

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: r,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				logger.Infof("HTTP server starting on port %d", cfg.Port)
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					logger.Errorf("HTTP server error: %v", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("HTTP server shutting down")
			return srv.Shutdown(ctx)
		},
	})

	return r
}
