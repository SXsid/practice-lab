package err

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type CreateUserRequest struct {
	Email    string `json:"email"  validate:"required,email"`
	UserName string `json:"user_name" validate:"required,min=5,max=30"`
}

func (app *App) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteErr(w, ErrIvalidFields)
		return
	}
	if err := app.validator.Struct(req); err != nil {
		WriteErr(w, err)
		return
	}
	// INFO:
	// passs the reqq.context as if use discoone this wil be  trigger the context cancel internally
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*3)
	defer cancel()
	if err := app.userservice.CreateUser(ctx, req.Email, req.UserName); err != nil {
		WriteErr(w, err)
		return
	}
	writeOk(w, http.StatusCreated, "user created sucessfully")
}
