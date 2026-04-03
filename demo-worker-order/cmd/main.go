package main

import (
	"flag"
	"fmt"

	"go.uber.org/fx"

	"github.com/gohive/core/logger"
	"github.com/gohive/demo-worker-order/config"
	"github.com/gohive/demo-worker-order/fxmodule"
)

func main() {
	configPath := flag.String("config", "./demo-worker-order/config.toml", "path to config file")
	flag.Parse()

	if err := config.Load(*configPath); err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}

	cfg := config.Cfg
	logger.Init(cfg.Log.Level, cfg.Log.Format)
	logger.Infof("Starting %s...", cfg.App.Name)

	app := fx.New(
		fx.Provide(
			fxmodule.NewDatabase,
			fxmodule.NewRedis,
			fxmodule.NewTaskRunner,
		),
		fx.Invoke(fxmodule.StartWorker),
	)

	app.Run()
}
