package err

import (
	"encoding/json"
	"net/http"
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
	// here too to match the naem is a pattern used or not
	if err := app.userservice.CreateUser(req.Email, req.UserName); err != nil {
		WriteErr(w, err)
		return
	}
	WriteResp(w, http.StatusCreated, "user created sucessfully", nil)
}
