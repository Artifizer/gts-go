/*
Copyright © 2025 Global Type System
Released under Apache License 2.0
*/

package server

import (
	"log"
	"net/http"
	"time"
)

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// withLogging wraps the handler with request logging
func (s *Server) withLogging(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if s.verbose == 0 {
			handler.ServeHTTP(w, r)
			return
		}

		start := time.Now()
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		handler.ServeHTTP(wrapped, r)

		duration := time.Since(start)
		log.Printf("%s %s -> %d in %.1fms",
			r.Method,
			r.URL.Path,
			wrapped.statusCode,
			float64(duration.Microseconds())/1000.0,
		)
	})
}
