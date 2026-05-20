package main

import (
	"fmt"
	"net/http"
)

var USERS = map[string]User{}

func main() {
	app := NewApp()
	server := http.Server{
		Addr:    ":8080",
		Handler: app.router,
	}
	fmt.Println("sever is up and running")
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
