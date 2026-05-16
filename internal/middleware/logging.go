package middleware

import (
	"context"
	"log"
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

		log.Printf(
			"[REQ] %s %s | IP: %s | RequestID: %s",
			r.Method,
			r.URL.Path,
			r.RemoteAddr,
			requestID,
		)

		next.ServeHTTP(rw, r)

		log.Printf(
			"[RES] %s %s | Status: %d | Duration: %v | RequestID: %s",
			r.Method,
			r.URL.Path,
			rw.StatusCode,
			time.Since(start),
			requestID,
		)
	})
}
