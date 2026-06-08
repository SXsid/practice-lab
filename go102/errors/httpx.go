package err

import (
	"encoding/json"
	"log"
	"net/http"
)

type ApiResponse struct {
	Success bool          `json:"success"`
	Message string        `json:"message,omitempty"`
	Body    any           `json:"body,omitempty"`
	Error   string        `json:"error,omitempty"`
	Field   []FieldErrors `json:"fields,omitempty"`
}

func write(w http.ResponseWriter, code int, body any) {
	b, _ := json.Marshal(body)
	w.Header().Set("Server", "go")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(b)
}

func WriteErr(w http.ResponseWriter, err error) {
	log.Printf("[error]:\n\t%w", err.Error())
	code, usrMessage := Resolve(err)
	field := ExtractFields(err)
	write(w, code, ApiResponse{
		Success: false,
		Error:   usrMessage,
		Field:   field,
	})
}

func writeOk(w http.ResponseWriter, code int, msg string, body ...any) {
	var data any
	if len(body) > 0 {
		data = body[0]
	}

	write(w, code, ApiResponse{
		Success: true,
		Message: msg,
		Body:    data,
	})
}
