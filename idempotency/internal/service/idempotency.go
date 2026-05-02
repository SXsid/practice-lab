package service

import (
	"context"
	"time"

	"github/SXsid/learn-idempotency/internal/domain"
)

type IdempotencyService interface {
	Get(ctx context.Context, key string) (*domain.IdempotencyRecord, error)
	Claim(ctx context.Context, key string, requestHash string, ttl time.Duration) (bool, error)
	Finalise(ctx context.Context, key string, response []byte, statusCode int) error
	Delete(ctx context.Context, key string) error
}
