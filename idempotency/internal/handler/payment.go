package handler

import (
	"context"
	"net/http"

	"github/SXsid/learn-idempotency/internal/domain"
)

type PaymentService interface {
	InitPayment(ctx context.Context, idempotencyID, requestHash string, customerID domain.CustomerID, amount int64, Currency domain.Currency) (string, error)
	HandleWebHook(ctx context.Context, OrderId string) error
	InitRefund(ctx context.Context, paymentID string) error
	ListPaymentByUser(ctx context.Context, customerID domain.CustomerID) ([]*domain.Payment, error)
}

type PaymentHandler struct {
	payService PaymentService
}

func NewPaymentHandler(payService PaymentService) *PaymentHandler {
	return &PaymentHandler{
		payService: payService,
	}
}

func (h *PaymentHandler) InitPayment(w http.ResponseWriter, r *http.Request) {
}

func (h *PaymentHandler) ProcessWebHook(w http.ResponseWriter, r *http.Request) {
}

func (h *PaymentHandler) InitiateRefund(w http.ResponseWriter, r *http.Request) {}
func (h *PaymentHandler) AllPayment(w http.ResponseWriter, r *http.Request)     {}
