package main

import (
	"sync"
	"time"
)

// INFO:  why?
// caues we have alot of work which doing sequite will kill performace or not optmial
// and we can't do all pearely / sequeltiy cause of limit and other extenal facotre
// so gate keep the goruinte execution called semphore

func main() {
	input := 1000
	var wg sync.WaitGroup
	// firin n go ruoint of n input
	// which could be bad
	// to limit / gate keep we use semaphore
	// 5 rouine concurent
	sem := make(chan struct{}, 5)
	for range input {
		wg.Go(func() {
			// find the slot
			sem <- struct{}{}
			// do the work
			time.Sleep(time.Millisecond * 100)
			<-sem
		})
	}
	wg.Wait()
	close(sem)
}
