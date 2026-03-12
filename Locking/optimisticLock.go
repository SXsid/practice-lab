package main

import (
	"fmt"
	"sync"
	"time"
)

type OptimisticWallet struct {
	Balance int
	version int
	mu      *sync.Mutex
}

func NewOptimisticWallet() *OptimisticWallet {
	return &OptimisticWallet{
		Balance: 100,
		version: 1,
		mu:      &sync.Mutex{},
	}
}

func (w *OptimisticWallet) Withdraw(amount int) {
	// read
	w.mu.Lock()
	balance := w.Balance
	version := w.version
	fmt.Println(balance)
	fmt.Println(version)
	w.mu.Unlock()

	// perfrm the operation
	if balance < amount {
		fmt.Println("Not enough funds")
	}
	newBalance := balance - amount
	time.Sleep(100 * time.Millisecond)
	// write
	w.mu.Lock()
	fmt.Println(w.version)
	fmt.Println(newBalance)
	if version != w.version {
		fmt.Println("data go updated while transction were on going")
		w.mu.Unlock()
		return
	}
	w.Balance = newBalance
	w.version++
	w.mu.Unlock()
}

func RunOptimisticBalance() {
	fmt.Println("start")
	wallet := NewOptimisticWallet()
	var wg sync.WaitGroup
	go func() {
		defer wg.Done()

		wallet.Withdraw(50)
	}()
	go func() {
		defer wg.Done()
		wallet.Withdraw(50)
	}()
	wg.Add(2)
	wg.Wait()
}
