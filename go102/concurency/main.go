package main

import (
	"fmt"
	"sync"
)

type User struct {
	Name string
}
type Result struct {
	user *User
	Err  error
}

func fetchUser(id int, ch chan<- Result) {
	if id%2 == 0 {
		ch <- Result{
			Err: fmt.Errorf("not a vaid id mate"),
		}
		return
	}
	ch <- Result{
		user: &User{Name: fmt.Sprintf("user:%d", id)},
	}
}

func main() {
	ch := make(chan Result)
	// sender on the other go rouinte  if it will be on same we cant read
	go func() {
		var wg sync.WaitGroup
		wg.Add(3)
		go func() { defer wg.Done(); fetchUser(1, ch) }()
		go func() { defer wg.Done(); fetchUser(2, ch) }()
		go func() { defer wg.Done(); fetchUser(3, ch) }()
		wg.Wait()
		// sende will tell we not gonan send more data
		close(ch)
	}()
	for res := range ch {
		user := res.user
		err := res.Err
		if err != nil {
			fmt.Println(err)
			continue

		}
		fmt.Println(user.Name)

	}
}
