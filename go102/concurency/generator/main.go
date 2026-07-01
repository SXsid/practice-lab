package main

import "fmt"

func generator[T any](arr []T) <-chan T {
	ch := make(chan T)
	// infite for loop
	go func() {
		for _, value := range arr {
			// stuck till we read kind of like yield
			ch <- value
		}
	}()
	return ch
}

func main() {
	name := []string{"sid", "harsh", "nirmal"}
	ch := generator(name)
	fmt.Println(<-ch)
	fmt.Println(<-ch)
	fmt.Println(<-ch)
}
