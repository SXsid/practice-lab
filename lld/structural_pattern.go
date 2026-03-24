package main

import (
	"errors"
	"fmt"
)

// ============================================================
// SINGLETON — Config loaded once, passed around
// ============================================================

type Config struct {
	Env      string // "local" | "prod" | "test"
	Host     string
	Port     int
	SSLMode  bool
	PoolSize int
}

var instance *Config // package level, created once

func NewConfig(env string) *Config {
	if instance != nil {
		return instance // already exists, return same
	}
	// INFO: in idea we will read form the env
	switch env {
	case "prod":
		instance = &Config{
			Env:      "prod",
			Host:     "prod.rds.amazonaws.com",
			Port:     5432,
			SSLMode:  true,
			PoolSize: 100,
		}
	case "test":
		instance = &Config{
			Env:      "test",
			Host:     "localhost",
			Port:     5432,
			SSLMode:  false,
			PoolSize: 1,
		}
	default: // local
		instance = &Config{
			Env:      "local",
			Host:     "localhost",
			Port:     5432,
			SSLMode:  false,
			PoolSize: 5,
		}
	}
	return instance
}

// ============================================================
// BUILDER — Knows HOW to construct a DB connection string
// ============================================================

type DatabaseBuilder struct {
	host     string
	port     int
	ssl      bool
	poolSize int
}

func NewDatabaseBuilder() *DatabaseBuilder {
	return &DatabaseBuilder{
		port:     5432, // sensible defaults
		poolSize: 10,
	}
}

func (b *DatabaseBuilder) WithHost(host string) *DatabaseBuilder {
	b.host = host
	return b
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

func (b *DatabaseBuilder) Build() (*Database, error) {
	if b.host == "" {
		return nil, errors.New("host is required") // validates before constructing
	}
	connStr := fmt.Sprintf(
		"postgres://%s:%d?sslmode=%v&pool_max_conns=%d",
		b.host, b.port, b.ssl, b.poolSize,
	)
	return &Database{connStr: connStr}, nil
}

// ============================================================
// DATABASE — the object we're building
// ============================================================

type Database struct {
	connStr string
}

func (d *Database) String() string {
	return fmt.Sprintf("Database connected: %s", d.connStr)
}

// ============================================================
// FACTORY — Decides WHICH config to use, delegates to Builder
// ============================================================
// Factory doesn't know HOW to build — that's Builder's job
// Builder doesn't know WHICH config to pick — that's Factory's job
// Clean separation of concerns

func NewDatabase(cfg *Config) (*Database, error) {
	// Factory decides which values to use based on env
	// Builder handles the actual construction
	return NewDatabaseBuilder().
		WithHost(cfg.Host).
		WithPort(cfg.Port).
		WithSSL(cfg.SSLMode).
		WithPoolSize(cfg.PoolSize).
		Build()
}

// ============================================================
// MAIN — Composition root, wires everything together
// ============================================================

func main() {
	// Singleton — created once, passed around
	cfg := NewConfig("prod")
	fmt.Printf("Config loaded for env: %s\n", cfg.Env)

	// verify singleton — same instance returned
	cfg2 := NewConfig("local")                     // tries to create local, but...
	fmt.Printf("Same instance? %v\n", cfg == cfg2) // true — prod config returned

	// Factory + Builder — factory picks config, builder constructs
	db, err := NewDatabase(cfg)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println(db)

	// Builder used directly — explicit, readable, validated
	customDB, err := NewDatabaseBuilder().
		WithHost("custom.host.com").
		WithSSL(true).
		WithPoolSize(50).
		Build()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println(customDB)
}
