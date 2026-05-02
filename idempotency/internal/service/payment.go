package service

import (
	"context"
	"time"

	"github/SXsid/learn-idempotency/internal/domain"

	"github.com/google/uuid"
)

type PaymentRepo interface {
	GetByID(ctx context.Context, ID string) (*domain.Payment, error)
	GetByOrderID(ctx context.Context, orderID domain.OrderID) (*domain.Payment, error)
	Create(ctx context.Context, payment domain.Payment) (*domain.Payment, error)
	UpdateProviderID(ctx context.Context, ID string, ProviderChargeID string, UpdatedAt time.Time) error
	UpdateStatus(ctx context.Context, ID string, status domain.PaymentStatus, UpdatedAt time.Time) error
	ListPayment(ctx context.Context, customerID domain.CustomerID) ([]*domain.Payment, error)
}

// this are method tha in servie we need not what razorpay or mock wiil ahve like we wiillimpem tem as a wroapper of htis funtin
type PaymentProvider interface {
	CreateOrder(ctx context.Context, Currency domain.Currency, amount int64) (string, error)
	Refund(ctx context.Context, orderId string, ProviderId string) error
}

type PaymentStore interface {
	Get(ctx context.Context) (*domain.PaymentIdempotency, error)
	Claim(ctx context.Context) error
	Finalise(ctx context.Context) error
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
	orderId, err := s.provider.CreateOrder(ctx, Currency, amount)
	if err != nil {
		return "", err
	}
	now := time.Now()
	paymentStruct := domain.Payment{
		ID:         uuid.New().String(),
		Amount:     amount,
		CustomerID: customerID,
		Currency:   Currency,
		Status:     domain.Pending,
		OrderID:    domain.OrderID(orderId),
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if _, err := s.repo.Create(ctx, paymentStruct); err != nil {
		return "", err
	}
	return orderId, nil
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
