package main

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	client *redis.Client
}

func NewRedisStore(ctx context.Context, addr, password string) (*RedisStore, error) {
	// Prod config options can be added here
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0, // use default DB
		PoolSize: 10,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	log.Println("Connected to Redis")
	return &RedisStore{client: client}, nil
}

// chain data use hset inseted of set for storing json
// if it's a lbog stick with set
func (s *RedisStore) Close() error {
	return s.client.Close()
}
