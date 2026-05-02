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

type PaymentService struct {
	repo     PaymentRepo
	provider PaymentProvider
	idemp    IdempotencyService
}

func NewPaymentService(repo PaymentRepo, provider PaymentProvider, idemp IdempotencyService) *PaymentService {
	return &PaymentService{
		repo:     repo,
		provider: provider,
		idemp:    idemp,
	}
}

func (s *PaymentService) InitPayment(ctx context.Context, idempotencyID, requestHash string, customerID domain.CustomerID, amount int64, Currency domain.Currency) (string, error) {
	// for nwo
	idemRecord, err := s.idemp.Get(ctx, idempotencyID)
	if idemRecord != nil {
		return string(idemRecord.Response), nil
	}
	claimed, err := s.idemp.Claim(ctx, idempotencyID, requestHash, time.Minute*10)
	if err != nil {
		return "", err
	}
	if !claimed {
		return "", domain.ErrRequestInFlight
	}

	orderId, err := s.provider.CreateOrder(ctx, Currency, amount)
	if err != nil {
		s.idemp.Delete(ctx, idempotencyID)
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
	if _, err = s.repo.Create(ctx, paymentStruct); err != nil {
		return "", err
	}
	if err := s.idemp.Finalise(ctx, idempotencyID, []byte(orderId), 201); err != nil {
		s.idemp.Delete(ctx, idempotencyID)
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
