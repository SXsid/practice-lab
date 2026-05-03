package dto

type InitPaymentResponse struct {
	OrderID string `json:"order_id"`
}

type InitPaymentRequest struct {
	CustomerID string  `json:"customer_id"`
	Currency   string  `json:"currency"`
	Amount     float64 `json:"amount"`
}
