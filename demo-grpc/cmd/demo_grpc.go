package cmd

import (
	"fmt"

	"go.uber.org/fx"
	"google.golang.org/grpc"
	"gorm.io/gorm"

	"github.com/gohive/core/logger"
	"github.com/gohive/demo-grpc/config"
	"github.com/gohive/demo-grpc/fxmodule"
	"github.com/gohive/demo-grpc/service"
)

func Run(configPath string) {
	// Load config
	if err := config.Load(configPath); err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}

	// Init logger
	cfg := config.Cfg
	logger.Init(cfg.Log.Level, cfg.Log.Format)

	// Create and run fx application
	app := fx.New(
		fx.Provide(
			fxmodule.NewDatabase,
			fxmodule.NewGRPCServer,
			service.NewUserService,
		),
		fx.Invoke(registerServices),
	)

	app.Run()
}

func registerServices(server *grpc.Server, db *gorm.DB, userService *service.UserService) {
	userService.Register(server)
	logger.Info("gRPC services registered")
}
