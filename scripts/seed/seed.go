package seed

import (
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/gohive/core/config"
	"github.com/gohive/core/logger"
	"github.com/gohive/models/entity"
	dbmysql "github.com/gohive/pkg/mysql"
)

// Options for seed command
type Options struct {
	ConfigPath string
	Force      bool // Force re-seed even if data exists
}

// Run executes database seeding.
func Run(opts Options) error {
	cfg, err := config.LoadConfig[config.BaseConfig](opts.ConfigPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	logger.Init(cfg.Log.Level, cfg.Log.Format)

	db, err := dbmysql.OpenGorm(cfg.Database)
	if err != nil {
		return fmt.Errorf("failed to connect database: %w", err)
	}
	defer dbmysql.CloseGorm(db)

	return runSeeder(db, opts.Force)
}

func runSeeder(db *gorm.DB, force bool) error {
	logger.Info("Running database seeder...")

	if err := seedUsers(db, force); err != nil {
		return err
	}

	if err := seedOrders(db, force); err != nil {
		return err
	}

	logger.Info("Database seeding completed successfully")
	return nil
}

func seedUsers(db *gorm.DB, force bool) error {
	var count int64
	if err := db.Model(&entity.User{}).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to count users: %w", err)
	}

	if count > 0 && !force {
		logger.Info("Users table already has data, skipping (use --force to re-seed)")
		return nil
	}

	if force && count > 0 {
		logger.Warn("Force mode: clearing existing users")
		if err := db.Exec("TRUNCATE TABLE users").Error; err != nil {
			return fmt.Errorf("failed to truncate users: %w", err)
		}
	}

	users := []entity.User{
		{
			Username: "admin",
			Email:    "admin@example.com",
			Password: "$2a$10$XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX", // hashed
			Nickname: "Administrator",
			Status:   1,
		},
		{
			Username: "user1",
			Email:    "user1@example.com",
			Password: "$2a$10$XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
			Nickname: "User One",
			Status:   1,
		},
		{
			Username: "user2",
			Email:    "user2@example.com",
			Password: "$2a$10$XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
			Nickname: "User Two",
			Status:   1,
		},
	}

	for _, user := range users {
		if err := db.Create(&user).Error; err != nil {
			return fmt.Errorf("failed to seed user %s: %w", user.Username, err)
		}
		logger.Infof("Seeded user: %s", user.Username)
	}

	return nil
}

func seedOrders(db *gorm.DB, force bool) error {
	var count int64
	if err := db.Model(&entity.Order{}).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to count orders: %w", err)
	}

	if count > 0 && !force {
		logger.Info("Orders table already has data, skipping (use --force to re-seed)")
		return nil
	}

	if force && count > 0 {
		logger.Warn("Force mode: clearing existing orders")
		if err := db.Exec("TRUNCATE TABLE orders").Error; err != nil {
			return fmt.Errorf("failed to truncate orders: %w", err)
		}
	}

	orders := []entity.Order{
		{
			OrderNo:   "ORD-20240101-001",
			UserID:    1,
			Amount:    99.99,
			Status:    1,
			ExpiredAt: time.Now().Add(24 * time.Hour),
		},
		{
			OrderNo:   "ORD-20240101-002",
			UserID:    2,
			Amount:    199.99,
			Status:    0,
			ExpiredAt: time.Now().Add(24 * time.Hour),
		},
	}

	for _, order := range orders {
		if err := db.Create(&order).Error; err != nil {
			return fmt.Errorf("failed to seed order %s: %w", order.OrderNo, err)
		}
		logger.Infof("Seeded order: %s", order.OrderNo)
	}

	return nil
}
