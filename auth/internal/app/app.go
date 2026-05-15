package app

import (
	"context"
)

type Application struct {
	Logger *Logger
}

func NewApp() (*Application, error) {
	ctx := context.Background()
	config := NewConfig()
	logger := NewLogger()
	db, err := NewPostgres(ctx, config.DSN)
	if err != nil {
		return nil, err
	}

	if err := Test(ctx, db); err != nil {
		return nil, err
	}
	return &Application{
		Logger: logger,
	}, nil
}

func (app *Application) Close() {
}
