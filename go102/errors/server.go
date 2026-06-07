package err

import "net/http"

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
}

func NewApp() *App {
	router := http.NewServeMux()
	app := &App{
		router:      router,
		userservice: NewUserService(NewUserRepo()),
	}
	router.HandleFunc("POST /createuser", app.CreateUser)
	return app
}
