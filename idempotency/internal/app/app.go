package app

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Application struct {
	config *Config
	logger *Logger
	db     *pgxpool.Pool
	rdc    *redis.Client
}

func NewApplicaton() http.Handler {
	ctx := context.Background()
	config := NewConfig()
	logger := NewLogger(config.LogLevel)
	db, err := NewPgxPool(ctx, config.DSN)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("postgress connecteed.")
	rdc, err := NewRedisClient(ctx, config.RedisURL)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("redis connecteed.")
	app := &Application{
		config: config,
		logger: logger,
		db:     db,
		rdc:    rdc,
	}
	return app.NewRouter()
}
