package app

import (
	"net/http"
)

type Application struct {
	config *Config
	logger *Logger
}

func NewApplicaton() http.Handler {
	config := NewConfig()
	logger := NewLogger(config.LogLevel)
	app := &Application{
		config: config,
		logger: logger,
	}
	return app.NewRouter()
}
