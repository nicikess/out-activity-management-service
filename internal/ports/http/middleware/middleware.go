package middleware

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

// Logger middleware for HTTP request logging
func Logger(log *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := &responseWriter{w: w, status: http.StatusOK}

			next.ServeHTTP(ww, r)

			log.Info("HTTP Request",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Int("status", ww.status),
				zap.Duration("duration", time.Since(start)),
				zap.String("remote_addr", r.RemoteAddr),
			)
		})
	}
}

// Auth middleware for authentication
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: Implement proper authentication
		// This is a placeholder - you should implement proper JWT validation
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		// For now, we'll just pass through the X-User-ID header
		// In a real implementation, you would validate the JWT and extract the user ID
		if r.Header.Get("X-User-ID") == "" {
			http.Error(w, "user ID not found", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// responseWriter is a wrapper around http.ResponseWriter that captures the status code
type responseWriter struct {
	w      http.ResponseWriter
	status int
}

func (rw *responseWriter) Header() http.Header {
	return rw.w.Header()
}

func (rw *responseWriter) Write(bytes []byte) (int, error) {
	return rw.w.Write(bytes)
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.status = statusCode
	rw.w.WriteHeader(statusCode)
}
