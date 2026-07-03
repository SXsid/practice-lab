package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	server := http.NewServeMux()
	server.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		select {
		case <-time.After(time.Second * 2):
			w.Write([]byte("<h1>Hello ladies</h1>"))
		case <-ctx.Done():
			fmt.Println("request cancelled")

		}
		fmt.Println("handler returned")
	})
	http.ListenAndServe(":8080", server)
}
