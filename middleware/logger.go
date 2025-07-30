package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

// StructuredLogger mencatat detail setiap request masuk.
func StructuredLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// Log request masuk
		slog.Info("Request diterima",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("remote_addr", r.RemoteAddr),
		)

		next.ServeHTTP(w, r)

		// Log request selesai
		slog.Info("Request selesai",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Duration("duration", time.Since(start)),
		)
	})
}

