package middlewares

import (
	"bytes"
	"net/http"
)

type customWriter struct {
	w      http.ResponseWriter
	data   bytes.Buffer
	status int
}

func NewCustomWrite(w http.ResponseWriter) *customWriter {
	return &customWriter{
		w:      w,
		status: 200,
	}
}

func (c *customWriter) Write(data []byte) (int, error) {
	c.data.Write(data)
	return c.w.Write(data)
}

func (c *customWriter) WriteHeader(status int) {
	c.status = status
	c.w.WriteHeader(status)
}

func (c *customWriter) Header() http.Header {
	return c.w.Header()
}

func IdempotencyMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cw := NewCustomWrite(w)
			next.ServeHTTP(cw, r)
		})
	}
}
