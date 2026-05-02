package provider

import (
	"context"
	"sync"

	"github/SXsid/learn-idempotency/internal/domain"

	"github.com/google/uuid"
)

type MockProvider struct {
	mu        sync.Mutex
	callCount int
}

func NewMockPayProvider() *MockProvider {
	return &MockProvider{
		callCount: 0,
	}
}

func (p *MockProvider) CreateOrder(ctx context.Context, Currency domain.Currency, amount int64) (string, error) {
	p.mu.Lock()
	p.callCount++
	p.mu.Unlock()
	uuid, err := uuid.NewRandom()
	if err != nil {
		return "", domain.ErrProviderDown
	}
	return uuid.String(), nil
}

func (p *MockProvider) Refund(ctx context.Context, orderId string, ProviderId string) error {
	return nil
}

func (p *MockProvider) CallCount() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.callCount
}
