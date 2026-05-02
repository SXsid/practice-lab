package domain

import "time"

type IdempotencyRecord struct {
	ID          string
	RequestHash string
	StausCode   int
	Response    []byte
	InFlight    bool
	CreatedAt   time.Time
	ExpiresAt   time.Time
}
