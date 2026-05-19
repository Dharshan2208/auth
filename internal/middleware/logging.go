package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		requestID := uuid.New().String()[:8]

		ctx := context.WithValue(
			r.Context(),
			"requestID",
			requestID,
		)

		r = r.WithContext(ctx)

		w.Header().Set("X-Request-ID", requestID)

		rw := &ResponseWriter{
			ResponseWriter: w,
		}

		next.ServeHTTP(rw, r)

		slog.Info("request ",
			"method", r.Method,
			"path", r.URL.Path,
			"status", rw.StatusCode,
			"duration", time.Since(start),
			"request_id", requestID,
		)
	})
}
