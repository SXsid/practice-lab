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

	"github/SXsid/learn-idempotency/internal/app"
)

func main() {
	port := flag.Int("port", 8080, "port at which app is running")
	flag.Parse()
	application := app.NewApplicaton()
	defer application.Close()
	server := http.Server{
		Addr:         fmt.Sprintf(":%d", *port),
		IdleTimeout:  1 * time.Minute,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  5 * time.Second,
		Handler:      app.NewRouter(application),
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
