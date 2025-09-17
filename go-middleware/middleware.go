package main

import (
	"fmt"
	"net/http"
	"os"
	"runtime/debug"
)

func RecoverFromPanicMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				stackTract := debug.Stack()
				fmt.Printf("Panic Recoverd: %v \n stack Trace: %s", r, stackTract)
				fmt.Println("key", os.Getenv("ENV"))
				if os.Getenv("ENV") == "development" {
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Fprintf(w, "Panic Reoverd ,\n Stack Trace:%s", stackTract)
					return
				}
				http.Error(w, "something went wrong", http.StatusInternalServerError)
			}
		}()
		nrw := &responseWriter{
			ResponseWriter: w,
			StatusCode:     http.StatusOK,
			Written:        false,
		}
		next.ServeHTTP(nrw, r)
		// if never panicked
		if err := nrw.flush(); err != nil {
			fmt.Printf("error couured while write %v", err)
		}
	}
}
