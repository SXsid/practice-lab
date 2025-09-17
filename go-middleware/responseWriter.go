package main

import "net/http"

type responseWriter struct {
	http.ResponseWriter
	StatusCode int
	Written    bool
	writes     [][]byte
}

func (rw *responseWriter) WriteHeader(statuscode int) {
	// we can't set header after we wrote the response
	if !rw.Written {
		rw.StatusCode = statuscode
	}
}

func (rw *responseWriter) Write(p []byte) (int, error) {
	if !rw.Written {
		rw.Written = true
	}
	rw.writes = append(rw.writes, p)
	return len(p), nil
}

func (rw *responseWriter) flush() error {
	rw.ResponseWriter.WriteHeader(rw.StatusCode)
	for _, data := range rw.writes {
		if _, err := rw.ResponseWriter.Write(data); err != nil {
			return err
		}
	}
	return nil
}
