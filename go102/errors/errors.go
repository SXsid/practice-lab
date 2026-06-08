package err

import (
	"errors"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

var (
	ErrIvalidFields     = errors.New("ivalid inputs")
	ErrEmailAlredyExist = errors.New("email is already in use")
	ErrServer           = errors.New("server error occured")
	ErrNotFound         = errors.New("not found")
)

type FieldErrors struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func Resolve(err error) (int, string) {
	switch {
	case errors.Is(err, ErrIvalidFields):
		return http.StatusUnprocessableEntity, ErrIvalidFields.Error()
	case errors.Is(err, ErrEmailAlredyExist):
		return http.StatusConflict, ErrEmailAlredyExist.Error()
	case errors.As(err, &validator.ValidationErrors{}):
		return http.StatusUnprocessableEntity, ErrIvalidFields.Error()
	default:
		return http.StatusInternalServerError, ErrServer.Error()
	}
}

func ExtractFields(err error) []FieldErrors {
	var ve validator.ValidationErrors
	if !errors.As(err, &ve) {
		return nil
	}
	fields := make([]FieldErrors, len(ve))
	for i, v := range ve {
		fields[i] = FieldErrors{
			Field:   strings.ToLower(v.Field()),
			Message: getMessage(v),
		}
	}

	return fields
}

func getMessage(v validator.FieldError) string {
	// deciide the msg with violation of which filed on erro at a tiem
	// for one field
	switch v.Tag() {
	case "required":
		return v.Field() + " is required"
	case "email":
		return v.Field() + " must be a valid email"
	case "min":
		return v.Field() + " must be atlest" + v.Param() + " characters"
	case "max":
		return v.Field() + " must be at most" + v.Param() + " charcters"

	default:
		return v.Field() + " is invalid"
	}
}
