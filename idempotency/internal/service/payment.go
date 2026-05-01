package service

import (
	"context"

	"github/SXsid/learn-idempotency/internal/domain"
)

type PaymentRepo interface {
	GetById(ctx context.Context, ID string) (*domain.Payment, error)
	GetByOrderId(ctx context.Context, orderID domain.OrderID) (*domain.Payment, error)
	Create(ctx context.Context, payment domain.Payment) (*domain.Payment, error)
	Update(ctx context.Context, payment domain.Payment) error
	ListPayment(ctx context.Context, customerID domain.CustomerID) ([]*domain.Payment, error)
}

// this are method tha in servie we need not what razorpay or mock wiil ahve like we wiillimpem tem as a wroapper of htis funtin
type PaymentProvider interface {
	CreateOrder() (string, error)
	Refund(orderId string, ProviderId string) error
}
