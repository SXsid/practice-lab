package err

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ErrResponse struct {
	Error string `json:"error"`
}

type Response struct {
	Message string `json:"msg"`
	Body    any    `json:"body"`
}

// here i gues sthe funtio optional pttern can be used
func WriteResp(w http.ResponseWriter, statusCode int, msg string, body any) {
	resp := Response{
		Message: msg,
		Body:    body,
	}
	b, _ := json.Marshal(resp)
	w.WriteHeader(statusCode)
	w.Write(b)
}

func WriteErr(w http.ResponseWriter, err error) {
	// log the acutl error
	fmt.Println(err.Error())
	code, msg := Resolve(err)
	resp := ErrResponse{
		Error: msg,
	}
	b, _ := json.Marshal(resp)
	w.WriteHeader(code)
	w.Write(b)
}
