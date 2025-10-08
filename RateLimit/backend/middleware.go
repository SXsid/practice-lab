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
		"/verfiy":         rate.NewLimiter(5, 5),
	}
)

func rateLimitMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL.Path)
		limter, ok := endPointLimiter[r.URL.Path]
		fmt.Println(ok)
		if ok {
			fmt.Println("inside url specifc")
			if !limter.Allow() {
				fmt.Println("inside this")
				http.Error(w, "Too Many Request", http.StatusTooManyRequests)
				return
			}
		} else {
			fmt.Println("inside global")
			if !globalLimiter.Allow() {
				fmt.Println("too many resps")
				http.Error(w, "Too Many Request", http.StatusTooManyRequests)
				return
			}
		}
		next.ServeHTTP(w, r)
	}
}
