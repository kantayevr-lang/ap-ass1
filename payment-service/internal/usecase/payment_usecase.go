package usecase

import (
	"context"
	"payment-service/internal/domain"

	"github.com/google/uuid"
)

type PaymentUseCase struct {
	repo domain.PaymentRepository
}

func NewPaymentUseCase(r domain.PaymentRepository) *PaymentUseCase {
	return &PaymentUseCase{repo: r}
}

func (u *PaymentUseCase) ProcessPayment(ctx context.Context, orderID string, amount int64) (*domain.Payment, error) {
	payment := &domain.Payment{
		ID:      uuid.New().String(),
		OrderID: orderID,
		Amount:  amount,
	}

	if amount > 100000 {
		payment.Status = domain.StatusDeclined
		payment.TransactionID = ""
	} else {
		payment.Status = domain.StatusAuthorized
		payment.TransactionID = uuid.New().String()
	}

	if err := u.repo.Save(ctx, payment); err != nil {
		return nil, err
	}

	return payment, nil
}
