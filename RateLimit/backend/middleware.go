package main

import (
	"fmt"
	"net/http"

	"golang.org/x/time/rate"
)

var (
	globalLimiter   = rate.NewLimiter(100, 170)
	endPointLimiter = map[string]*rate.Limiter{
		"/forgotPassword": rate.NewLimiter(10, 15),
		"/verify":         rate.NewLimiter(5, 5),
	}
)

func rateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !globalLimiter.Allow() {
			fmt.Println("blocking the request for ", r.RemoteAddr)
			http.Error(w, "Too Many Request", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func endPointRateLimiter(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limter, ok := endPointLimiter[r.URL.Path]
		fmt.Printf("inside the endpone thing")
		if ok {
			if !limter.Allow() {
				http.Error(w, "Too Many Request", http.StatusTooManyRequests)
				return
			}
		}
		next.ServeHTTP(w, r)
	}
}
