package main

import (
	"sync"
	"testing"
)

func TestDeposit(t *testing.T) {
	var wg sync.WaitGroup
	account := NewAccount(10_000, "1")
	for range 50 {
		wg.Add(1)
		go func() {
			account.Deposit(100)
			defer wg.Done()
		}()
	}
	wg.Wait()
	if account.amount != 15_000 {
		t.Errorf("expected 15,000 got %.2f", account.amount)
	}
}

func TestWithdraw(t *testing.T) {
	var wg sync.WaitGroup
	account := NewAccount(10_000, "1")
	for range 50 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			account.Withdraw(100)
		}()
	}
	wg.Wait()
	if account.Balance() != 5_000 {
		t.Errorf("expected 5,000 got %.2f", account.amount)
	}
}
