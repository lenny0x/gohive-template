package fxmodule

import (
	"context"
	"fmt"

	"go.uber.org/fx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/gohive/core/logger"
	"github.com/gohive/demo-api/config"
)

func NewDatabase(lc fx.Lifecycle) (*gorm.DB, error) {
	cfg := config.Cfg.Database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("Database connected")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			sqlDB, err := db.DB()
			if err != nil {
				return err
			}
			logger.Info("Database connection closed")
			return sqlDB.Close()
		},
	})

	return db, nil
}
