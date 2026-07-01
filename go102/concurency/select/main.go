package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type User struct {
	id int
}
type Result struct {
	u   *User
	err error
}

func fetchUser(ctx context.Context, id int, ch chan<- Result) {
	select {
	case <-time.After(time.Millisecond * 300):
		ch <- Result{&User{id: id}, nil}
	case <-ctx.Done():
		ch <- Result{nil, fmt.Errorf("time out")}
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*300)
	defer cancel()
	// INFO:
	// so we sender doesnt have to wait/stuck even wehn we close the reciver when erro get
	ch := make(chan Result, 3)
	go func() {
		var wg sync.WaitGroup
		for i := range 3 {
			wg.Go(func() {
				fetchUser(ctx, i, ch)
			})
		}
		wg.Wait()
		close(ch)
	}()
	for res := range ch {
		u := res.u
		if err := res.err; err != nil {
			fmt.Println(err)
			break
		}
		fmt.Println(u.id)

	}
}
