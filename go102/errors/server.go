package err

import (
	"net/http"

	"github.com/go-playground/validator/v10"
)

func Run() {
	app := NewApp()
	server := http.Server{
		Addr:    ":8080",
		Handler: app.router,
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

// TODO: https://www.youtube.com/watch?v=J1PDCaJrQG8
// watch that work on error then respo and rest cycle
// so couple of question soem coding the real interview perp
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
