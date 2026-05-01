package service

import (
	"context"
	"time"

	"github/SXsid/learn-idempotency/internal/domain"
)

type PaymentRepo interface {
	GetById(ctx context.Context, ID string) (*domain.Payment, error)
	GetByOrderId(ctx context.Context, orderID domain.OrderID) (*domain.Payment, error)
	Create(ctx context.Context, payment domain.Payment) (*domain.Payment, error)
	UpdateProviderId(ctx context.Context, ID string, ProviderChargeID string, UpdatedAt time.Time) error
	UpdateStatus(ctx context.Context, ID string, status domain.PaymentStatus, UpdatedAt time.Time) error
	ListPayment(ctx context.Context, customerID domain.CustomerID) ([]*domain.Payment, error)
}

// this are method tha in servie we need not what razorpay or mock wiil ahve like we wiillimpem tem as a wroapper of htis funtin
type PaymentProvider interface {
	CreateOrder(ctx context.Context, paymentID string, amount int64) (string, error)
	Refund(ctx context.Context, orderId string, ProviderId string) error
}

type PaymentService struct {
	repo     PaymentRepo
	provider PaymentProvider
}

func NewPaymentService(repo PaymentRepo, provider PaymentProvider) *PaymentService {
	return &PaymentService{
		repo:     repo,
		provider: provider,
	}
}

func (s *PaymentService) InitPayment(ctx context.Context, customerID domain.CustomerID, amount int64, Currency domain.Currency) (string, error) {
	return "", nil
}

func (s *PaymentService) HandleWebHook(ctx context.Context, OrderId string) error {
	return nil
}

func (s *PaymentService) InitRefund(ctx context.Context, paymentID string) error {
	return nil
}

func (s *PaymentService) ListPaymentByUser(ctx context.Context, customerID domain.CustomerID) ([]*domain.Payment, error) {
	return nil, nil
}
