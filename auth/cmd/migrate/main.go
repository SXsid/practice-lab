package main

import (
	"database/sql"
	"fmt"
	"os"

	"github/SXsid/auth-learn/migrations"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func main() {
	dsn := os.Getenv("DSN")
	if dsn == "" {
		fmt.Println("no dsn key exsist to migrate")
		os.Exit(1)
	}
	db, err := sql.Open("pgx", dsn)
	if err != nil {

		fmt.Println(err.Error())
		os.Exit(1)
	}
	goose.SetBaseFS(migrations.MigrationFS)
	defer goose.SetBaseFS(nil)
	if err := goose.SetDialect("postgres"); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if err := goose.Up(db, "."); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)

	}
}
