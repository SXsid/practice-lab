package app

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPgxPool(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		fmt.Println("dsn string is worgn")
		return nil, err
	}
	pgx, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		fmt.Println("postgress is not setup")
		return nil, err
	}
	if err := pgx.Ping(ctx); err != nil {
		fmt.Println("connection is not established")
		return nil, err
	}
	return pgx, nil
}
