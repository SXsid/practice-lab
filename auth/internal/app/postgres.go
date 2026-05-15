package app

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgres(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}
	return pgxpool.NewWithConfig(ctx, cfg)
}

func Test(ctx context.Context, pgx *pgxpool.Pool) error {
	ctxTime, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	return pgx.Ping(ctxTime)
}
