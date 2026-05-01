package handler

import (
	"context"

	"github/SXsid/learn-idempotency/internal/domain"
)

type PaymentSerive interface {
	InitPayment(ctx context.Context, amount int64) (string, error)

	// ther will be the paymeprovide struct what data it send iidk abu ti
	HandleWebHook(ctx context.Context, OrderId string)
	InitRefund(ctx context.Context, paymentID string) error
	ListPaymentByUser(ctx context.Context, customerID domain.CustomerID)
}
