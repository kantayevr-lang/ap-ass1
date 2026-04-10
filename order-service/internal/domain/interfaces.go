package domain

import "context"

type OrderRepository interface {
	Save(ctx context.Context, order *Order) error
	GetByID(ctx context.Context, id string) (*Order, error)
	UpdateStatus(ctx context.Context, id string, status string) error
}

type PaymentService interface {
	Pay(ctx context.Context, orderID string, amount int64) (string, error)
}
