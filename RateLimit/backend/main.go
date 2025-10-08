package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"net/http"
)

var otp = make(map[string]string)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", handleHome)
	mux.HandleFunc("GET /forgotPassword", forgorPassword)
	mux.HandleFunc("GET /verify", VerifyOtp)
	fmt.Println("server is up at prot:8080")
	if err := http.ListenAndServe(":8080", rateLimitMiddleware(mux)); err != nil {
		log.Fatalf("error while starting the server")
	}
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	_, _ = fmt.Fprintln(w, "<h1>Hello World!!</h1>")
}

func forgorPassword(w http.ResponseWriter, r *http.Request) {
	var code string
	for i := 0; i < 4; i++ {

		r, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			fmt.Printf("err occured%v", err.Error())
			i--
			continue
		}
		code += fmt.Sprint(r)
	}
	fmt.Printf("code:%s\n", code)
	otp["8003592767"] = code
}

func VerifyOtp(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("otp")
	if code == otp["8003592767"] {
		fmt.Println("verified")
		w.WriteHeader(http.StatusOK)
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}
