package usecase

import (
	"context"
	"order-service/internal/domain"
)

type OrderUseCase struct {
	repo          domain.OrderRepository
	paymentClient domain.PaymentService
}

func NewOrderUseCase(r domain.OrderRepository, p domain.PaymentService) *OrderUseCase {
	return &OrderUseCase{
		repo:          r,
		paymentClient: p,
	}
}

func (u *OrderUseCase) CreateOrder(ctx context.Context, ord *domain.Order) error {
	if ord.Amount <= 0 {
		return domain.ErrInvalidAmount
	}

	ord.Status = domain.StatusPending

	if err := u.repo.Save(ctx, ord); err != nil {
		return err
	}

	status, err := u.paymentClient.Pay(ctx, ord.ID, ord.Amount)
	if err != nil {
		ord.Status = domain.StatusFailed
		u.repo.UpdateStatus(ctx, ord.ID, ord.Status)
		return err
	}

	finalStatus := domain.StatusFailed
	if status == "Authorized" {
		finalStatus = domain.StatusPaid
	}

	ord.Status = finalStatus
	return u.repo.UpdateStatus(ctx, ord.ID, ord.Status)
}

func (u *OrderUseCase) CancelOrder(ctx context.Context, id string) error {
	order, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if order.Status != domain.StatusPending {
		return domain.ErrCannotCancelOrder
	}

	return u.repo.UpdateStatus(ctx, id, domain.StatusCancelled)
}
