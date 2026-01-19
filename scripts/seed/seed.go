package seed

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/gohive/core/config"
	"github.com/gohive/core/logger"
	"github.com/gohive/models/entity"
)

// Options for seed command
type Options struct {
	ConfigPath string
	Force      bool // Force re-seed even if data exists
}

// Run executes database seeding
func Run(opts Options) {
	cfg, err := config.LoadConfig[config.BaseConfig](opts.ConfigPath)
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

	// Run seeder
	if err := runSeeder(db, opts.Force); err != nil {
		logger.Fatalf("Seeding failed: %v", err)
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
		db.Exec("TRUNCATE TABLE users")
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
			logger.Errorf("Failed to seed user %s: %v", user.Username, err)
			return err
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
		db.Exec("TRUNCATE TABLE orders")
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
			logger.Errorf("Failed to seed order %s: %v", order.OrderNo, err)
			return err
		}
		logger.Infof("Seeded order: %s", order.OrderNo)
	}

	return nil
}
