package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type App struct {
	Router *chi.Mux
	pg     *PostgresStore
	redis  *RedisStore
	config *Config
	logger *Logger
}

func NewApp(cfg *Config, pg *PostgresStore, rds *RedisStore, logger *Logger) *App {
	app := &App{
		Router: chi.NewRouter(),
		pg:     pg,
		redis:  rds,
		config: cfg,
		logger: logger,
	}

	app.setupRoutes()
	return app
}

func (a *App) setupRoutes() {
	a.Router.Use(middleware.RequestID)
	a.Router.Use(middleware.RealIP)
	a.Router.Use(middleware.Logger)
	a.Router.Use(middleware.Recoverer)
	a.Router.Use(middleware.Timeout(60 * time.Second))

	a.Router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// Example Leaderboard routes
	a.Router.Post("/leaderboard", a.handleUpdateScore)
	a.Router.Get("/leaderboard", a.handleGetLeaderboard)
}

func (a *App) handleUpdateScore(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (a *App) handleGetLeaderboard(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
