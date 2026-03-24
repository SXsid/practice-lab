package main

import (
	"errors"
	"fmt"
	"time"
)

type DatabaseBuilder struct {
	host       string
	port       int
	ssl        bool
	poolSize   int
	timeout    time.Duration
	maxRetries int
}

func NewDatabase() *DatabaseBuilder {
	return &DatabaseBuilder{
		port:       5432, // sensible defaults
		poolSize:   10,
		timeout:    30 * time.Second,
		maxRetries: 3,
	}
}

// each method sets one field and returns the builder
func (b *DatabaseBuilder) WithHost(host string) *DatabaseBuilder {
	b.host = host
	return b // INFO: returns self for chaining
}

func (b *DatabaseBuilder) WithPort(port int) *DatabaseBuilder {
	b.port = port
	return b
}

func (b *DatabaseBuilder) WithSSL(ssl bool) *DatabaseBuilder {
	b.ssl = ssl
	return b
}

func (b *DatabaseBuilder) WithPoolSize(size int) *DatabaseBuilder {
	b.poolSize = size
	return b
}

// Build validates and constructs the final object
func (b *DatabaseBuilder) Build() (*pgxpool.Pool, error) {
	if b.host == "" {
		return nil, errors.New("host is required")
	}
	connStr := fmt.Sprintf(
		"postgres://%s:%d?sslmode=%v&pool_max_conns=%d",
		b.host, b.port, b.ssl, b.poolSize,
	)
	return pgxpool.New(context.Background(), connStr)
}
