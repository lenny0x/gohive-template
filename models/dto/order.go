package dto

type CreateOrderRequest struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
}

type OrderResponse struct {
	ID        uint    `json:"id"`
	OrderNo   string  `json:"order_no"`
	UserID    uint    `json:"user_id"`
	Amount    float64 `json:"amount"`
	Status    int     `json:"status"`
	CreatedAt string  `json:"created_at"`
}
