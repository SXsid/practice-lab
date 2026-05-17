package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// LeaderboardEntry represents a struct for live leaderboard app boilerplate
type LeaderboardEntry struct {
	UserID   string  `json:"user_id"`
	Username string  `json:"username"`
	Score    float64 `json:"score"`
	Rank     int64   `json:"rank,omitempty"`
}

func main() {
	ctx := context.Background()

	// Load configuration
	cfg := LoadConfig()
	logger := NewLogger()

	// Initialize Postgres
	// You might want to handle an empty DSN gracefully in a real app
	pgStore, err := NewPostgresStore(ctx, cfg.PostgresDSN)
	if err != nil {
		log.Printf("Failed to initialize Postgres (Make sure POSTGRES_DSN is set correctly): %v", err)
	} else {
		defer pgStore.Close()
	}

	// Initialize Redis
	redisStore, err := NewRedisStore(ctx, cfg.RedisAddr, cfg.RedisPass)
	if err != nil {
		log.Printf("Failed to initialize Redis (Make sure REDIS_ADDR is set correctly): %v", err)
	} else {
		defer redisStore.Close()
	}

	// Initialize App
	app := NewApp(cfg, pgStore, redisStore, logger)

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: app.Router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("error while starting app..", "error", err.Error())

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
