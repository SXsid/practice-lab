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
	ch := make(chan Result)
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
			continue
		}
		fmt.Println(u.id)

	}
}
