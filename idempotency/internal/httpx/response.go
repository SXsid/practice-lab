package httpx

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Response struct {
	Error  string         `json:"error,omitempty"`
	Data   any            `json:"data,omitempty"`
	Msg    string         `json:"msg,omitempty"`
	Fields map[string]any `json:"fields,omitempty"`
}

func WriteError(w http.ResponseWriter, errMsg string, status int, fields map[string]any) {
	res := Response{
		Error:  errMsg,
		Fields: fields,
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Server", "Go")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		fmt.Print(err)
	}
}

func WriteResponse(w http.ResponseWriter, msg string, status int, data any) {
	res := Response{
		Data: data,
		Msg:  msg,
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Server", "Go")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		fmt.Print(err)
	}
}
