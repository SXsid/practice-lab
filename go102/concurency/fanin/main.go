package main

import (
	"fmt"
	"sync"
)

// INFO: IDEA
// multiple  prodcuer
// we act as consmer
func generator[T any](arr []T) <-chan T {
	ch := make(chan T)
	go func() {
		defer close(ch)
		for _, value := range arr {
			// hanged here
			ch <- value
		}
	}()
	return ch
}

func main() {
	// two producer read the input form them process and reutn the ruslt
	i1 := generator([]int{1, 2, 3, 4})
	i2 := generator([]int{5, 6, 7, 8, 9, 0})
	res := fanIn(i1, i2)
	for value := range res {
		fmt.Println("value:", value)
	}
}

func fanIn(inputs ...<-chan int) <-chan int {
	out := make(chan int)
	// make consumer two rouine with unbuufer chanenl
	// so two  consuerm
	// both consumer don't know when it will end
	// so routine can't cose the channel
	// we hvae to wait when both / all consuemr rouitne done
	// we close the ruin
	// co ordinator  to wait on multipel consumer
	go func() {
		var wg sync.WaitGroup
		for _, input := range inputs {
			wg.Add(1)
			go func(ch <-chan int) {
				defer wg.Done()
				for value := range ch {
					out <- value * 10
				}
			}(input)
		}
		wg.Wait()
		close(out)
	}()
	return out
}
