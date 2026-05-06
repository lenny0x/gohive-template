package fixorder

import (
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/gohive/core/config"
	"github.com/gohive/core/logger"
	"github.com/gohive/models/entity"
	dbmysql "github.com/gohive/pkg/mysql"
)

// Order status constants for the orders.status column.
const (
	statusPending   = 0
	statusCancelled = 3
)

// Options for fix-order command
type Options struct {
	ConfigPath string
	DryRun     bool // Only show what would be changed
}

// Run executes the order fix script.
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

	return fixExpiredOrders(db, opts.DryRun)
}

func fixExpiredOrders(db *gorm.DB, dryRun bool) error {
	if dryRun {
		logger.Warn("DRY RUN MODE - No changes will be made")
	}

	logger.Info("Finding expired pending orders...")

	// Pin "now" once so the SELECT and UPDATE see the same cutoff and
	// can't race against orders crossing the expiry boundary mid-run.
	now := time.Now()

	var orders []entity.Order
	err := db.Where("status = ? AND expired_at < ?", statusPending, now).Find(&orders).Error
	if err != nil {
		return fmt.Errorf("failed to query orders: %w", err)
	}

	if len(orders) == 0 {
		logger.Info("No expired pending orders found")
		return nil
	}

	logger.Infof("Found %d expired pending orders", len(orders))

	ids := make([]uint, 0, len(orders))
	for _, order := range orders {
		ids = append(ids, order.ID)
		logger.Infof("  - Order %s (User: %d, Amount: %.2f, Expired: %s)",
			order.OrderNo, order.UserID, order.Amount, order.ExpiredAt.Format(time.RFC3339))
	}

	if dryRun {
		logger.Warnf("DRY RUN: Would update %d orders to cancelled status", len(orders))
		return nil
	}

	// Update only the rows we just selected, by id, to avoid TOCTOU drift
	// against orders crossing the expiry boundary between SELECT and UPDATE.
	result := db.Model(&entity.Order{}).
		Where("id IN ? AND status = ?", ids, statusPending).
		Update("status", statusCancelled)

	if result.Error != nil {
		return fmt.Errorf("failed to update orders: %w", result.Error)
	}

	logger.Infof("Successfully updated %d orders to cancelled status", result.RowsAffected)
	return nil
}
