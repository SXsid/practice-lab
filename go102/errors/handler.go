package err

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type CreateUserRequest struct {
	Email    string `json:"email"`
	UserName string `json:"user_name"`
}

func (c *CreateUserRequest) Validate() error {
	var ers []error
	if c.Email == "" || !strings.Contains(c.Email, "@") {
		ers = append(ers, fmt.Errorf("email is not valid or empty"))
	}
	if len(c.UserName) < 5 {
		ers = append(ers, fmt.Errorf("len of username can't be less than 5"))
	}
	if len(ers) > 0 {
		res := []error{
			ErrIvalidFields,
		}

		res = append(res, ers...)
		return errors.Join(res...)
	}
	return nil
}

func (app *App) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteErr(w, ErrIvalidFields)
		return
	}
	if err := req.Validate(); err != nil {
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
