package domain

import "time"

type (
	CustomerID string
	OrderID    string
	Currency   string
)

type PaymentStatus string

const (
	Pending PaymentStatus = "pending"
	Success PaymentStatus = "suceeded"
	Failed  PaymentStatus = "failed"
)

type PaymentIdemPotencyStatus string

const (
	Miss     PaymentIdemPotencyStatus = "miss"
	InFlight PaymentIdemPotencyStatus = "inFlight"
	Hit      PaymentIdemPotencyStatus = "hit"
)

type Payment struct {
	ID               string
	CustomerID       CustomerID
	OrderID          OrderID
	Amount           int64
	Currency         Currency
	ProviderChargeID *string
	Status           PaymentStatus
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type PaymentIdempotency struct {
	Id     string
	status PaymentIdemPotencyStatus
}

func (p *Payment) GetAmount() float64 {
	return float64(p.Amount) / float64(100)
}

func (p *Payment) MarkSuccess() error {
	if p.Status != Pending {
		return ErrInvalidStatusTransition
	}
	p.Status = Success
	return nil
}

func (p *Payment) MarkFailed() error {
	if p.Status != Pending {
		return ErrInvalidStatusTransition
	}
	p.Status = Failed
	return nil
}
