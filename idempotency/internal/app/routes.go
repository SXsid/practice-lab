package app

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func NewRouter(app *Application) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"OPTIONS", "GET", "POST"},
		AllowedHeaders: []string{"Origin", "Accept", "Content-Type"},
		// for dev in 0
		MaxAge: int((time.Hour * 0).Seconds()),
	}))

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("helo word 2.0"))
	})
	router.Mount("/api/v1", V1Routes(app))
	return router
}

func V1Routes(app *Application) http.Handler {
	r := chi.NewRouter()
	r.Mount("/payment", PayRotues(app))
	return r
}

func PayRotues(app *Application) http.Handler {
	r := chi.NewRouter()
	r.Post("/init", app.payHandler.InitPayment)
	r.Patch("/webhook", app.payHandler.ProcessWebHook)
	r.Post("/refund", app.payHandler.InitiateRefund)
	r.Get("/all", app.payHandler.AllPayment)
	return r
}
