package provider

import (
	"context"
	"sync"

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

func (p *MockProvider) CreateOrder(ctx context.Context, paymentID string, amount int64) (string, error) {
	p.mu.Lock()
	p.callCount++
	p.mu.Unlock()
	uuid, err := uuid.NewRandom()
	return uuid.String(), err
}

func (p *MockProvider) Refund(ctx context.Context, orderId string, ProviderId string) error {
	return nil
}

func (p *MockProvider) CallCount() int {
	return p.callCount
}
