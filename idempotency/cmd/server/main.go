package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("helo word"))
	})
	port := flag.Int("port", 8080, "port at which app is running")
	flag.Parse()
	server := http.Server{
		Addr:         fmt.Sprintf(":%d", *port),
		IdleTimeout:  1 * time.Minute,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  5 * time.Second,
		Handler:      router,
	}
	go func() {
		fmt.Printf("server is up and running at %s \n", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("err stating the server %s", err.Error())
			os.Exit(1)
		}
	}()
	sctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	<-sctx.Done()
	fmt.Println("server stopping gracefully")
	tctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	if err := server.Shutdown(tctx); err != nil {
		fmt.Println("err while stopping server gracefully")
		os.Exit(1)
	}
	fmt.Println("server stopped.")
}
