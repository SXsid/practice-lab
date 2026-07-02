package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// INFO: IDEA
// multiple conusmer
// we act as producer

func generator(ctx context.Context) <-chan int {
	ch := make(chan int)
	go func() {
		count := 0
		for {
			select {
			case <-ctx.Done():
				close(ch)
				return
			case ch <- count:
				count++

			}
		}
	}()
	return ch
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()
	producer := generator(ctx)
	c1 := consumer(producer)
	c2 := consumer(producer)
	c3 := consumer(producer)
	c4 := consumer(producer)
	// we can't use select as we don't know mych we have to read
	merged := FanIn(c1, c2, c3, c4)
	for value := range merged {
		fmt.Println("count:", value)
	}
}

// merge the chaneel to revve the output
// INFO: one chanenl for all go rouen ,so a corrinate closes after we are done readin all
func FanIn(inputs ...<-chan int) <-chan int {
	out := make(chan int)
	// unfferd array so as many  consumer for the consmer lol
	go func() {
		var wg sync.WaitGroup
		for _, input := range inputs {
			wg.Go(func() {
				// readin those parrel worker   parrel and when done  we move out of rouine
				for i := range input {
					out <- i
				}
			})
		}
		wg.Wait()
		// done reading
		close(out)
	}()
	return out
}

// INFO: invidi l consumer so invideo channel closes
func consumer(ch <-chan int) <-chan int {
	out := make(chan int)

	go func(ch <-chan int) {
		defer close(out)
		for value := range ch {
			// heavey work
			time.Sleep(100 * time.Millisecond)
			out <- value
		}
	}(ch)

	return out
}
