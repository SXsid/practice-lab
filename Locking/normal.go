package main

import (
	"fmt"
	"sync"
)

type NormalWallet struct {
	Balance int
}

func NewNormalWallet() *NormalWallet {
	return &NormalWallet{
		Balance: 100,
	}
}

func (w *NormalWallet) Withdraw(amount int) error {
	if w.Balance < amount {
		return fmt.Errorf("not enough fund")
	}
	w.Balance -= amount
	return nil
}

func RunNormal() {
	fmt.Println("start")
	wallet := NewNormalWallet()
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()

		fmt.Println(wallet.Balance)
		if err := wallet.Withdraw(50); err != nil {
			fmt.Println(err)
		}
		fmt.Println(wallet.Balance)
	}()
	go func() {
		defer wg.Done()
		fmt.Println(wallet.Balance)
		if err := wallet.Withdraw(60); err != nil {
			fmt.Println(err)
		}
		fmt.Println(wallet.Balance)
	}()
	wg.Wait()
}
