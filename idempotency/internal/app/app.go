package app

import (
	"context"
	"fmt"
	"os"

	"github/SXsid/learn-idempotency/internal/handler"
	"github/SXsid/learn-idempotency/internal/provider"
	"github/SXsid/learn-idempotency/internal/repository/postgres"
	"github/SXsid/learn-idempotency/internal/service"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Application struct {
	config     *Config
	logger     *Logger
	db         *pgxpool.Pool
	rdc        *redis.Client
	payHandler *handler.PaymentHandler
}

func NewApplicaton() *Application {
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
	payRepo := postgres.NewPaymentRepo(db)
	payProvider := provider.NewMockPayProvider()
	payService := service.NewPaymentService(payRepo, payProvider)
	payHandler := handler.NewPaymentHandler(payService)
	app := &Application{
		config:     config,
		logger:     logger,
		db:         db,
		rdc:        rdc,
		payHandler: payHandler,
	}
	return app
}

func (app *Application) Close() {
	app.db.Close()
	app.rdc.Close()
}
