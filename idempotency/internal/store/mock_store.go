package store

import (
	"context"
	"sync"
	"time"

	"github/SXsid/learn-idempotency/internal/domain"
)

type MockStore struct {
	mu   sync.Mutex
	data map[string]*domain.IdempotencyRecord
}

func NewMockStore() *MockStore {
	return &MockStore{
		data: make(map[string]*domain.IdempotencyRecord),
	}
}

func (s *MockStore) Get(ctx context.Context, key string) (*domain.IdempotencyRecord, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	val, ok := s.data[key]
	if !ok {
		return nil, nil
	}
	copy := *val
	return &copy, nil
}

func (s *MockStore) Claim(ctx context.Context, key string, requestHash string, ttl time.Duration) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.data[key]
	if ok {
		return false, nil
	}

	now := time.Now()
	data := domain.IdempotencyRecord{
		ID:          key,
		RequestHash: requestHash,
		InFlight:    true,
		CreatedAt:   now,
		ExpiresAt:   now.Add(ttl),
	}
	s.data[key] = &data
	return true, nil
}

func (s *MockStore) Finalise(ctx context.Context, key string, response []byte, statusCode int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	val, ok := s.data[key]
	if !ok {
		return domain.ErrNotFound
	}
	val.InFlight = false
	val.Response = response
	val.StausCode = statusCode
	return nil
}

func (s *MockStore) Delete(ctx context.Context, key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, key)
	return nil
}
