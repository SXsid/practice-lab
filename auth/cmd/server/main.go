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

	"github/SXsid/auth-learn/internal/api/router"
	"github/SXsid/auth-learn/internal/app"
)

func main() {
	port := flag.Int("port", 9090, "port of the http server")
	flag.Parse()

	app, err := app.NewApp()
	if err != nil {
		panic(err)
	}
	defer app.Close()
	logger := app.Logger

	server := http.Server{
		Addr:         fmt.Sprintf(":%d", *port),
		IdleTimeout:  1 * time.Minute,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  5 * time.Second,
		Handler:      router.NewRouter(app),
	}
	go func() {
		logger.Info("server is up and running at port ", "PORT", *port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println(err.Error())
			os.Exit(1)

		}
	}()

	sigCTX, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	<-sigCTX.Done()
	logger.Info("gracefully shutting down server..")
	timeCTX, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	if err := server.Shutdown(timeCTX); err != nil {
		logger.Error("error while gracful shutdown ", "ERROR", err)
		os.Exit(1)
	}
}
