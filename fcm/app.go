package main

import (
	"encoding/json"
	"net/http"

	firebase "firebase.google.com/go/v4"
)

type App struct {
	router   *http.ServeMux
	User     map[string]User
	firebase *firebase.App
}

func NewApp() *App {
	router := http.NewServeMux()
	app := &App{
		router: router,
		User: map[string]User{
			"602b6dc3-72b0-4916-9336-aea7d1cdf39c": {
				ID:    "602b6dc3-72b0-4916-9336-aea7d1cdf39c",
				Name:  "Harsh",
				Token: "cwpAhrRhmtX370xM2iinV5:APA91bHsDvyG8YeB3n8k0kFQ7poY6aYP-geiK8zlUjMtup-p7YQoSdRO3az5gkP-4E-JxoqD4-AjSZl_hMt5CNr_6ejjfmWxX5jCjr0ukfb27hsY8ckdfXk",
			}, "a7dbffc3-8768-4287-9566-799d0f72dde2": {
				ID:    "a7dbffc3-8768-4287-9566-799d0f72dde2",
				Name:  "Nirmal",
				Token: "cfE27-_D2nNygI57OLNUOq:APA91bHAVoFlHE-ptDnlo54iF7tu6rOG6DQE46YWJjG_yPzZT56I4H28m-O1AksdnDMtULXnNp_aXlo_tPgdVIO0eFHYEh-IUfBQH-8wot88Tv1K3hliGvA",
			},
		},
		firebase: NewFirebaseClient(),
	}
	app.SetupRoutes()
	return app
}

func (app *App) SetupRoutes() {
	fileHandler := http.FileServerFS(FinalStaticFS)
	app.router.Handle("GET /", fileHandler)
	app.router.HandleFunc("POST /api/v1/fcmToken", func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")
		ID := r.URL.Query().Get("id")
		user, ok := app.User[ID]
		if !ok {
			http.Error(w, "id not found", http.StatusNotFound)
			return
		}
		user.Token = token
		app.User[ID] = user
	})
	app.router.HandleFunc("GET /api/v1/profile/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		user, ok := app.User[id]
		if !ok {
			http.Error(w, ErrNotFound.Error(), http.StatusNotFound)
			return
		}
		data, err := json.Marshal(user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	})
	app.router.HandleFunc("POST /api/send", app.SendNotificaion)
}
