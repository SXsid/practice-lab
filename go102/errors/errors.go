package err

import (
	"errors"
	"net/http"
)

var (
	ErrIvalidFields     = errors.New("ivalid field")
	ErrEmailAlredyExist = errors.New("email is already in use")
	ErrServer           = errors.New("server error occured")
	ErrNotFound         = errors.New("not found")
)

func Resolve(err error) int {
	switch {
	case errors.Is(err, ErrIvalidFields), errors.Is(err, ErrEmailAlredyExist):
		return http.StatusUnprocessableEntity
	case errors.Is(err, ErrServer):
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}
