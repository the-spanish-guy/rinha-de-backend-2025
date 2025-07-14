package middleware

import (
	"net/http"
	"rinha-de-backend-2025/internal/config"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func Logging(logger *config.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			wrapper := &responseWriter{ResponseWriter: w, statusCode: 200}

			next.ServeHTTP(wrapper, r)

			duration := time.Since(start)
			logger.Infof("method: %s path: %s queryparams %s statusCode: %d duration(ms): %v",
				r.Method,
				r.URL.Path,
				r.URL.Query(),
				wrapper.statusCode,
				duration)

		})
	}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
