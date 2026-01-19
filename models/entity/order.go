package entity

import "time"

type Order struct {
	BaseModel
	OrderNo   string    `gorm:"size:64;uniqueIndex" json:"order_no"`
	UserID    uint      `gorm:"index" json:"user_id"`
	Amount    float64   `gorm:"type:decimal(10,2)" json:"amount"`
	Status    int       `gorm:"default:0" json:"status"` // 0: pending, 1: paid, 2: completed, 3: cancelled
	ExpiredAt time.Time `json:"expired_at"`
}

func (Order) TableName() string {
	return "orders"
}
