package postgres

import (
	"context"
	"time"

	"github/SXsid/learn-idempotency/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PaymentRepo struct {
	pgx *pgxpool.Pool
}

func NewPaymentRepo(pgx *pgxpool.Pool) *PaymentRepo {
	return &PaymentRepo{
		pgx: pgx,
	}
}

func (r *PaymentRepo) GetByID(ctx context.Context, ID string) (*domain.Payment, error) {
	return nil, nil
}

func (r *PaymentRepo) GetByOrderID(ctx context.Context, orderID domain.OrderID) (*domain.Payment, error) {
	return nil, nil
}

func (r *PaymentRepo) Create(ctx context.Context, payment domain.Payment) (*domain.Payment, error) {
	return nil, nil
}

func (r *PaymentRepo) UpdateProviderID(ctx context.Context, ID string, ProviderChargeID string, UpdatedAt time.Time) error {
	return nil
}

func (r *PaymentRepo) UpdateStatus(ctx context.Context, ID string, status domain.PaymentStatus, UpdatedAt time.Time) error {
	return nil
}

func (r *PaymentRepo) ListPayment(ctx context.Context, customerID domain.CustomerID) ([]*domain.Payment, error) {
	return nil, nil
}
