package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type DBPool struct {
	sem chan struct{} // your semaphore
}

func NewDBPool(maxConns int) *DBPool {
	return &DBPool{
		sem: make(chan struct{}, maxConns),
	}
}

func (p *DBPool) Acquire(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("timeout ,please try later")
	case p.sem <- struct{}{}:
		return nil
		// default:
		//
		// 	//INFO:  bad as i am wait  here forever  ctx is not respectd
		// 	// p.sem <- struct{}{}
		// 	return nil

	}
}

func (p *DBPool) Release() {
	// remove form buffer
	<-p.sem
}

func (p *DBPool) Avaiblabel() int {
	// INFO: the lent read is is atmoic so no need of lock
	return cap(p.sem) - len(p.sem)
}

func simulateDBQuery(id int, pool *DBPool, wg *sync.WaitGroup) {
	defer wg.Done()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := pool.Acquire(ctx); err != nil {
		fmt.Printf("worker %d: gave up waiting — %v\n", id, err)
		return
	}
	defer pool.Release()

	fmt.Printf("worker %d: acquired connection\n", id)
	time.Sleep(500 * time.Millisecond) // simulate DB work
	fmt.Printf("worker %d: done, releasing\n", id)
}

func session2Init() {
	pool := NewDBPool(3) // only 3 concurrent connections allowed
	var wg sync.WaitGroup

	// 10 workers, only 3 can be in DB at once
	for i := 1; i <= 10; i++ {
		wg.Add(1)
		go simulateDBQuery(i, pool, &wg)
	}

	wg.Wait()
	fmt.Println("all done")
}
