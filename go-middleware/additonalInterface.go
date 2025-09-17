package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
)

// NOTE:it's tricky as the inner respose write don't return a error so we can handle that and stop using this mehtod
func (rw *responseWriter) Flush() {
	flusher, ok := rw.ResponseWriter.(http.Flusher)
	if ok {
		flusher.Flush()
		return
	}
	return
}

func (rw *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	// so this chekc is for the inner respose wite as the custom autmiilly implemt the hijacker cause of thi
	// fuction
	hijack, ok := rw.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, fmt.Errorf("inner responseWriter dosen't support hijacker")
	}
	return hijack.Hijack()
}
