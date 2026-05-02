package provider

import (
	"context"

	"github/SXsid/learn-idempotency/internal/domain"
)

type Razorpay struct{}

func NewRazopayClient() *Razorpay {
	return &Razorpay{}
}

func (r *Razorpay) CreateOrder(ctx context.Context, Currency domain.Currency, amount int64) (string, error) {
	return "order_id_for_pay", nil
}

func (r *Razorpay) Refund(ctx context.Context, orderId string, ProviderId string) error {
	return nil
}
