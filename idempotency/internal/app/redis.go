package app

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(ctx context.Context, redisURL string) (*redis.Client, error) {
	opt, err := redis.ParseURL(redisURL)
	opt.PoolSize = 50 // depend on the througput
	opt.MaxRetries = 3
	opt.MinIdleConns = 10
	if err != nil {
		fmt.Println("cant' setup rdis")
		return nil, err
	}
	rdc := redis.NewClient(opt)
	if err := rdc.Ping(ctx).Err(); err != nil {
		fmt.Println("reids server is not responding ")
		return nil, err
	}
	return rdc, nil
}
