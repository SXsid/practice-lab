package main

import (
	"context"
	"testing"
	"time"
)

// Test + happy path           → does it work at all
// 1 acqurie and avilbity
func TestAcquirdLockWhenFree(t *testing.T) {
	pollSize := 5
	pool := NewDBPool(pollSize)
	// no cancellation for happy path
	ctx := context.Background()
	if err := pool.Acquire(ctx); err != nil || pool.Avaiblabel() != 4 {
		t.Errorf("expectd nil and 4 got : %v %d", err, pool.Avaiblabel())
	}
}

// 2) relese slot : if we relse a buffer will another grap it
func TestRelesingFreeSlot(t *testing.T) {
	pollSize := 1
	pool := NewDBPool(pollSize)
	ctx := context.Background()

	_ = pool.Acquire(ctx)
	pool.Release()
	if err := pool.Acquire(ctx); err != nil {
		t.Errorf("expected nil got , %v", err)
	}
}

// Test +  eexpected / requreid error stimualte them

func TestAcquirdFailWhileeWaitingFilledPool(t *testing.T) {
	// timeout context
	pollSize := 1
	pool := NewDBPool(pollSize)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	go func() {
		_ = pool.Acquire(ctx)
		defer pool.Release()
		time.Sleep(4 * time.Second)
	}()
	time.Sleep(time.Millisecond * 10)
	if err := pool.Acquire(ctx); err == nil {
		t.Errorf("expectd timeout error got , %v", err)
	}
}
