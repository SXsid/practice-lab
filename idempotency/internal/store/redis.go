package store

import (
	"context"
	"time"

	"github/SXsid/learn-idempotency/internal/domain"

	"github.com/redis/go-redis/v9"
)

type Store struct {
	rdc *redis.Client
}

func NewRedisStore(rdc *redis.Client) *Store {
	return &Store{
		rdc: rdc,
	}
}

func (s *Store) Get(ctx context.Context, key string) (domain.IdempotencyRecord, error) {
	return domain.IdempotencyRecord{}, nil
}

func (s *Store) Claim(ctx context.Context, key string, requestHash string, ttl time.Duration) (bool, error) {
	return false, nil
}

func (s *Store) Finalise(ctx context.Context, key string, response []byte, statusCode int) error {
	return nil
}

func (s *Store) Delete(ctx context.Context, key string) error {
	return nil
}
