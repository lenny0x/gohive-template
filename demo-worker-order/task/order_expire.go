package task

import (
	"context"

	"github.com/gohive/core/logger"
)

type OrderExpireTask struct{}

func NewOrderExpireTask() *OrderExpireTask {
	return &OrderExpireTask{}
}

func (t *OrderExpireTask) Name() string {
	return "order_expire"
}

func (t *OrderExpireTask) Run(ctx context.Context) error {
	logger.Info("running order expire task...")
	// TODO: implement order expiration logic
	// - Query expired orders from database
	// - Update order status
	// - Send notifications
	return nil
}
