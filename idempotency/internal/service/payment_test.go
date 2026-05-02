package service

import (
	"context"
	"sync"
	"testing"

	"github/SXsid/learn-idempotency/internal/domain"
	"github/SXsid/learn-idempotency/internal/provider"
	"github/SXsid/learn-idempotency/internal/repository/mock"
)

func TestInitPayment_ConcurrentSameRequest_ShouldOnlyChargeOnce(t *testing.T) {
	// arrange
	repo := mock.NewMockRepo()
	provider := provider.NewMockPayProvider()
	svc := NewPaymentService(repo, provider)

	const concurrency = 10
	var wg sync.WaitGroup
	wg.Add(concurrency)

	// act — 10 goroutines, same customer, same amount
	for range concurrency {
		go func() {
			defer wg.Done()
			_, _ = svc.InitPayment(
				context.Background(),
				domain.CustomerID("cust_123"),
				10050,
				domain.Currency("INR"),
			)
		}()
	}
	wg.Wait()

	// assert — this is what SHOULD be true in a correct system
	// both will FAIL right now — that's the point

	if provider.CallCount() != 1 {
		t.Errorf("provider called %d times, want 1 — double charging happening", provider.CallCount())
	}

	if repo.Count() != 1 {
		t.Errorf("repo created %d records, want 1 — duplicate records in DB", repo.Count())
	}
}
