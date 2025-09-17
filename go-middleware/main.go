package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "<h1>Hello go!!</h1>")
	})
	mux.HandleFunc("GET /panic/", DemoWichWilPanic)
	if err := http.ListenAndServe(":8080", RecoverFromPanicMiddleware(mux)); err != nil {
		log.Fatalf("error occured when starting the serverr%s", err.Error())
	}
}

func DemoWichWilPanic(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "<h3>gonna panic</h3>")
	panicFunction()
}

// method which panics
func panicFunction() {
	panic("ohh no")
}
