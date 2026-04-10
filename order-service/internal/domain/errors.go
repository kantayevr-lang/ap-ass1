package domain

import "errors"

var (
	ErrOrderNotFound      = errors.New("order not found")
	ErrInvalidAmount      = errors.New("amount must be greater than zero")
	ErrCannotCancelOrder  = errors.New("cannot cancel paid order")
	ErrPaymentServiceDown = errors.New("payment service unavailable")
)
