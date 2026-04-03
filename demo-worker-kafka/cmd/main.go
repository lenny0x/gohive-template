package main

import (
	"flag"
	"fmt"

	"go.uber.org/fx"

	"github.com/gohive/core/logger"
	"github.com/gohive/demo-worker-kafka/config"
	"github.com/gohive/demo-worker-kafka/fxmodule"
)

func main() {
	configPath := flag.String("config", "./demo-worker-kafka/config.toml", "path to config file")
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
		),
		fx.Invoke(fxmodule.NewConsumer),
	)

	app.Run()
}
