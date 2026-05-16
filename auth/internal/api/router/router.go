package router

import (
	"net/http"

	"github/SXsid/auth-learn/internal/api/middlewares"
	"github/SXsid/auth-learn/internal/app"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func NewRouter(app *app.Application) http.Handler {
	chi := chi.NewRouter()
	chi.Use(middleware.Recoverer)
	chi.Use(middlewares.Correlation)
	chi.Get("/", func(w http.ResponseWriter, r *http.Request) {
		logger := middlewares.LoggerWithContext(r.Context(), app.Logger)
		logger.Info("requs came")
		w.Write([]byte("<h1> HELLO AUTH JI </h1>"))
	})
	return chi
}
