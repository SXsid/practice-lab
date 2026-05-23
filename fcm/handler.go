package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"firebase.google.com/go/v4/messaging"
)

func (app *App) RemoveToken(ID string) error {
	user, ok := app.User[ID]
	if !ok {
		return ErrNotFound
	}
	user.Token = ""
	app.User[ID] = user
	return nil
}

func (app *App) SendNotificaion(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	firebaseClient, err := app.firebase.Messaging(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	res := make([]string, 0)
	for _, u := range app.User {
		msg := messaging.Message{
			Token: u.Token,
			Notification: &messaging.Notification{
				Title: "result are out",
				Body:  "check you result live of test concducon 22 april",
			},
		}
		Tres, err := firebaseClient.Send(ctx, &msg)
		if err != nil {
			switch err {
			// case firebaseClient.
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return

		}
		fmt.Println("sent")
		res = append(res, Tres)

	}

	w.WriteHeader(http.StatusCreated)
	data, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return

	}
	w.Write(data)
}
