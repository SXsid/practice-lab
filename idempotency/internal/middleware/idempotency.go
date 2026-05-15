package middleware

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"time"

	"github/SXsid/learn-idempotency/internal/domain"
	"github/SXsid/learn-idempotency/internal/service"
)

type responseCapture struct {
	http.ResponseWriter
	header     http.Header
	statusCode int
	body       bytes.Buffer
}

func NewresponseCapture(w http.ResponseWriter) *responseCapture {
	return &responseCapture{
		ResponseWriter: w,
		statusCode:     200,
	}
}

func (r *responseCapture) WriteHeader(code int) {
	r.statusCode = code
	r.ResponseWriter.WriteHeader(code)
}

func (r *responseCapture) Write(body []byte) (int, error) {
	r.body.Write(body)
	return r.ResponseWriter.Write(body)
}

func (r *responseCapture) Header() http.Header {
	return r.ResponseWriter.Header()
}

func hashRequest(r *http.Request) (string, error) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return "", err
	}
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	h := sha256.New()
	h.Write([]byte(r.Method))
	h.Write([]byte(r.URL.Path))
	h.Write(bodyBytes)

	return hex.EncodeToString(h.Sum(nil)), nil
}

func IdempotecyMiddleware(idem service.IdempotencyService) func(http.Handler) http.Handler {
	// chi inject this
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			idem_id := r.Header.Get("idempotechy-key")
			if idem_id == "" {
				http.Error(w, domain.ErrInvalidBody.Error(), http.StatusBadRequest)
				return
			}
			ctx := r.Context()
			requestHash, err := hashRequest(r)
			if err != nil {
				http.Error(w, domain.ErrServerSide.Error(), http.StatusInternalServerError)
				return
			}
			idemRecord, err := idem.Get(ctx, idem_id)
			if err != nil {
				http.Error(w, domain.ErrServerSide.Error(), http.StatusInternalServerError)
				return
			}
			if idemRecord != nil {
				if idemRecord.RequestHash != requestHash {
					http.Error(w, "idempotency key reused with different request", http.StatusUnprocessableEntity) // 422
					return
				}
				if idemRecord.InFlight {
					http.Error(w, domain.ErrRequestInFlight.Error(), http.StatusConflict)
					return
				}
				w.WriteHeader(idemRecord.StausCode)
				w.Write(idemRecord.Response)
				return
			}
			claimed, err := idem.Claim(ctx, idem_id, requestHash, time.Minute*10)
			if err != nil {
				http.Error(w, domain.ErrServerSide.Error(), http.StatusInternalServerError)
				return
			}
			if !claimed {
				// double click solved
				http.Error(w, domain.ErrRequestInFlight.Error(), http.StatusConflict)
				return
			}
			rc := NewresponseCapture(w)
			next.ServeHTTP(rc, r)
			// INFO: cause it saved in reponse writer and body is not
			rc.header = rc.Header().Clone()
			// INFO: in user error we delete so retry can happen || iif wrong we regisite the stept n outbox event and idempdcy will contiibue from there
			// for handler iit's new as it was failed prev try
			if rc.statusCode >= 400 {
				idem.Delete(ctx, idem_id)
				return
			}
			if err := idem.Finalise(ctx, idem_id, rc.body.Bytes(), rc.statusCode); err != nil {
				idem.Delete(ctx, idem_id)
				http.Error(w, domain.ErrServerSide.Error(), http.StatusInternalServerError)
				return
			}
		})
	}
}
