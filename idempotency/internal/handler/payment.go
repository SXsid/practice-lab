package handler

import (
	"context"
	"encoding/json"
	"math"
	"net/http"

	"github/SXsid/learn-idempotency/internal/domain"
	dto "github/SXsid/learn-idempotency/internal/handler/DTO"
	"github/SXsid/learn-idempotency/internal/httpx"
)

type PaymentService interface {
	InitPayment(ctx context.Context, customerID domain.CustomerID, amount int64, Currency domain.Currency) (string, error)
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
	ctx := r.Context()
	var body dto.InitPaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		httpx.WriteError(w, domain.ErrInvalidBody.Error(), http.StatusBadRequest, nil)
		return
	}
	orderId, err := h.payService.InitPayment(ctx, domain.CustomerID(body.CustomerID), int64(math.Round(body.Amount*100)), domain.Currency(body.Currency))
	if err != nil {
		httpx.WriteError(w, domain.ErrServerSide.Error(), http.StatusInternalServerError, nil)
		return
	}
	data := dto.InitPaymentResponse{
		OrderID: orderId,
	}
	httpx.WriteResponse(w, "payemnt initiated", http.StatusCreated, data)
}

func (h *PaymentHandler) ProcessWebHook(w http.ResponseWriter, r *http.Request) {
}

func (h *PaymentHandler) InitiateRefund(w http.ResponseWriter, r *http.Request) {}
func (h *PaymentHandler) AllPayment(w http.ResponseWriter, r *http.Request)     {}
