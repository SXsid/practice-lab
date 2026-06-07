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

func Resolve(err error) (int, string) {
	switch {
	case errors.Is(err, ErrIvalidFields):
		return http.StatusUnprocessableEntity, ErrIvalidFields.Error()
	case errors.Is(err, ErrEmailAlredyExist):
		return http.StatusUnprocessableEntity, ErrEmailAlredyExist.Error()
	default:
		return http.StatusInternalServerError, ErrServer.Error()
	}
}
