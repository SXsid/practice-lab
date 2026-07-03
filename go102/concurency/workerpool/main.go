package main

import (
	"fmt"
	"sync"
	"time"
)

const WORKER_POOL = 4

func main() {
	input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	inputCh := generator(input)
	out := make(chan int)
	go func() {
		var wg sync.WaitGroup
		for range WORKER_POOL {
			wg.Go(func() {
				Worker(inputCh, out)
			})
		}
		wg.Wait()
		close(out)
	}()
	for res := range out {
		fmt.Println(res)
	}
}

func generator(arr []int) <-chan int {
	ch := make(chan int)
	go func() {
		defer close(ch)
		for _, value := range arr {
			ch <- value
		}
	}()
	return ch
}

func Worker(inputch <-chan int, res chan<- int) {
	for input := range inputch {
		time.Sleep(1000 * time.Millisecond)
		res <- 2 * input
	}
}
