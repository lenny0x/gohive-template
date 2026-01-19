package fixorder

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/gohive/core/config"
	"github.com/gohive/core/logger"
	"github.com/gohive/models/entity"
)

// Options for fix-order command
type Options struct {
	ConfigPath string
	DryRun     bool // Only show what would be changed
}

// Run executes the order fix script
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

	// Fix expired orders
	if err := fixExpiredOrders(db, opts.DryRun); err != nil {
		logger.Fatalf("Fix failed: %v", err)
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

func fixExpiredOrders(db *gorm.DB, dryRun bool) error {
	if dryRun {
		logger.Warn("DRY RUN MODE - No changes will be made")
	}

	logger.Info("Finding expired pending orders...")

	// Find orders that are pending (status=0) and expired
	var orders []entity.Order
	err := db.Where("status = ? AND expired_at < ?", 0, time.Now()).Find(&orders).Error
	if err != nil {
		logger.Errorf("Failed to query orders: %v", err)
		return err
	}

	if len(orders) == 0 {
		logger.Info("No expired pending orders found")
		return nil
	}

	logger.Infof("Found %d expired pending orders", len(orders))

	for _, order := range orders {
		logger.Infof("  - Order %s (User: %d, Amount: %.2f, Expired: %s)",
			order.OrderNo, order.UserID, order.Amount, order.ExpiredAt.Format(time.RFC3339))
	}

	if dryRun {
		logger.Warnf("DRY RUN: Would update %d orders to cancelled status", len(orders))
		return nil
	}

	// Update orders to cancelled status (status=3)
	result := db.Model(&entity.Order{}).
		Where("status = ? AND expired_at < ?", 0, time.Now()).
		Update("status", 3)

	if result.Error != nil {
		logger.Errorf("Failed to update orders: %v", result.Error)
		return result.Error
	}

	logger.Infof("Successfully updated %d orders to cancelled status", result.RowsAffected)
	return nil
}
