package main

import (
	"database/sql"
	"fmt"
	"os"

	"github/SXsid/learn-idempotency/migrations"

	_ "github.com/jackc/pgx/v5"
	"github.com/pressly/goose/v3"
)

func main() {
	DSN := os.Getenv("DSN")
	if DSN == "" {
		panic("")
	}
	conn, err := sql.Open("pgx", DSN)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	goose.SetBaseFS(migrations.MigrationFs)
	defer goose.SetBaseFS(nil)

	if err := goose.SetDialect("postgres"); err != nil {
		fmt.Println(err)
		os.Exit(1)

	}

	if err := goose.Up(conn, "."); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
