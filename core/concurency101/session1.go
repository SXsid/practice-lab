package main

import (
	"fmt"
	"sync"
)

type Account struct {
	mu     sync.Mutex
	amount float64
	id     string
}

func NewAccount(amount float64, id string) *Account {
	return &Account{
		amount: amount,
		id:     id,
	}
}

func (a *Account) Balance() float64 {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.amount
}

func (a *Account) Deposit(amount float64) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.amount += amount
}

func (a *Account) Withdraw(amount float64) error {
	if amount < 0 {
		return fmt.Errorf("amount can't -ve")
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	userBalance := a.amount
	if userBalance < amount {
		return fmt.Errorf("insufficent funds")
	}
	a.amount = userBalance - amount
	return nil
}
