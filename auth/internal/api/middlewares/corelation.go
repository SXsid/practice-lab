package middlewares

import (
	"context"
	"net/http"

	"github/SXsid/auth-learn/internal/app"

	"github.com/google/uuid"
)

type corelationContextKey struct{}

func Correlation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		correlationID := uuid.NewString()
		// INFO: don't use the build iint data type define custom type to avoid collison with other libs
		// we use struct as it relally hard to be comon instead of type key stirng const colreco key =" conet"
		ctx := context.WithValue(r.Context(), corelationContextKey{}, correlationID)
		// new r  as it also imutanle as the the context
		r = r.WithContext(ctx)
		w.Header().Set("X-correlation-ID", correlationID)
		next.ServeHTTP(w, r)
	})
}

func GetCorreltaionId(ctx context.Context) string {
	id, ok := ctx.Value(corelationContextKey{}).(string)
	if !ok {
		return "Unknown"
	}
	return id
}

func LoggerWithContext(ctx context.Context, logger *app.Logger) *app.Logger {
	return &app.Logger{
		Logger: logger.Logger.With("corrleation_id", GetCorreltaionId(ctx)),
	}
}
