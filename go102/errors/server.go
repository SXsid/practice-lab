package err

import (
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
)

func Run() {
	app := NewApp()
	server := http.Server{
		Addr:         ":8080",
		Handler:      app.router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}

type App struct {
	userservice *UserService
	router      http.Handler
	validator   *validator.Validate
}

func NewApp() *App {
	router := http.NewServeMux()
	app := &App{
		router:      router,
		userservice: NewUserService(NewUserRepo()),
		validator:   validator.New(),
	}
	router.HandleFunc("POST /createuser", app.CreateUser)
	return app
}
