package mock

import (
	"context"
	"sync"
	"time"

	"github/SXsid/learn-idempotency/internal/domain"
)

type MockRepo struct {
	mu    sync.RWMutex
	data  map[string]*domain.Payment
	count int64
}

func NewMockRepo() *MockRepo {
	return &MockRepo{
		data: make(map[string]*domain.Payment),
	}
}

func (m *MockRepo) Count() int64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.count
}

func (m *MockRepo) Create(ctx context.Context, payment domain.Payment) (*domain.Payment, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	// INFO: we use orderID as more look up goona hapen usign this so o(1) and for id it's o(n)
	m.data[string(payment.OrderID)] = &payment
	m.count++
	return &payment, nil
}

func (m *MockRepo) GetByID(ctx context.Context, ID string) (*domain.Payment, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, p := range m.data {
		if p.ID == ID {
			return p, nil
		}
	}
	return nil, nil
}

func (m *MockRepo) GetByOrderID(ctx context.Context, orderID domain.OrderID) (*domain.Payment, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	p, ok := m.data[string(orderID)]
	if !ok {
		return nil, nil
	}
	return p, nil
}

func (m *MockRepo) UpdateStatus(ctx context.Context, ID string, status domain.PaymentStatus, updatedAt time.Time) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, p := range m.data {
		if p.ID == ID {
			p.Status = status
			p.UpdatedAt = updatedAt
			return nil
		}
	}
	return nil
}

func (m *MockRepo) UpdateProviderID(ctx context.Context, ID string, providerChargeID string, updatedAt time.Time) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, p := range m.data {
		if p.ID == ID {
			p.ProviderChargeID = &providerChargeID
			p.UpdatedAt = updatedAt
			return nil
		}
	}
	return nil
}

func (m *MockRepo) ListPayment(ctx context.Context, customerID domain.CustomerID) ([]*domain.Payment, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var result []*domain.Payment
	for _, p := range m.data {
		if p.CustomerID == customerID {
			result = append(result, p)
		}
	}
	return result, nil
}
