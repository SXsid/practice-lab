package main

import (
	"fmt"
	"sync"
)

type PessimitcWallet struct {
	Balance int
	mutext  *sync.Mutex
}

func NewPessimisticWallet() *PessimitcWallet {
	return &PessimitcWallet{
		Balance: 100,
		mutext:  &sync.Mutex{},
	}
}

func (w *PessimitcWallet) Withdraw(amount int) error {
	w.mutext.Lock()
	defer w.mutext.Unlock()
	fmt.Println(w.Balance)

	if w.Balance < amount {
		return fmt.Errorf("not enough fund")
	}

	w.Balance -= amount
	fmt.Println(w.Balance)

	return nil
}

func RunPesimistic() {
	fmt.Println("start")
	wallet := NewPessimisticWallet()
	var wg sync.WaitGroup
	go func() {
		defer wg.Done()

		if err := wallet.Withdraw(50); err != nil {
			fmt.Println(err)
		}
	}()
	go func() {
		defer wg.Done()
		if err := wallet.Withdraw(60); err != nil {
			fmt.Println(err)
		}
	}()
	wg.Add(2)
	wg.Wait()
}
