package migrate

import (
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"

	"github.com/gohive/core/config"
	"github.com/gohive/core/logger"
	dbmysql "github.com/gohive/pkg/mysql"
)

func setup(configPath string) (*sql.DB, error) {
	cfg, err := config.LoadConfig[config.BaseConfig](configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	logger.Init(cfg.Log.Level, cfg.Log.Format)

	if err := goose.SetDialect("mysql"); err != nil {
		return nil, fmt.Errorf("failed to set dialect: %w", err)
	}

	db, err := dbmysql.OpenSQL(cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	return db, nil
}

func Up(configPath, migrationsDir string) error {
	db, err := setup(configPath)
	if err != nil {
		return err
	}
	defer db.Close()
	return goose.Up(db, migrationsDir)
}

func Down(configPath, migrationsDir string) error {
	db, err := setup(configPath)
	if err != nil {
		return err
	}
	defer db.Close()
	return goose.Down(db, migrationsDir)
}

func Status(configPath, migrationsDir string) error {
	db, err := setup(configPath)
	if err != nil {
		return err
	}
	defer db.Close()
	return goose.Status(db, migrationsDir)
}

func Create(migrationsDir, name, migrationType string) error {
	return goose.Create(nil, migrationsDir, name, migrationType)
}
