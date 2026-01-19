package migrate

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/gohive/core/config"
	"github.com/gohive/core/logger"
	"github.com/gohive/models/entity"
)

// Run executes database migrations
func Run(configPath string) {
	cfg, err := config.LoadConfig[config.BaseConfig](configPath)
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}

	// Init logger
	logger.Init(cfg.Log.Level, cfg.Log.Format)

	// Connect database
	db, err := connectDB(cfg.Database)
	if err != nil {
		logger.Fatalf("Failed to connect database: %v", err)
	}
	defer closeDB(db)

	// Run migrations
	if err := runMigrations(db); err != nil {
		logger.Fatalf("Migration failed: %v", err)
	}
}

func connectDB(cfg config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}

func closeDB(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		return
	}
	sqlDB.Close()
}

func runMigrations(db *gorm.DB) error {
	logger.Info("Running database migrations...")

	models := []interface{}{
		&entity.User{},
		&entity.Order{},
	}

	for _, model := range models {
		if err := db.AutoMigrate(model); err != nil {
			logger.Errorf("Failed to migrate %T: %v", model, err)
			return err
		}
		logger.Infof("Migrated: %T", model)
	}

	logger.Info("Database migrations completed successfully")
	return nil
}
