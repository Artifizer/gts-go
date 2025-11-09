/*
Copyright © 2025 Global Type System
Released under Apache License 2.0
*/

package server

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	body       bytes.Buffer
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(p []byte) (int, error) {
	// Capture body
	rw.body.Write(p)
	return rw.ResponseWriter.Write(p)
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

		// If highest verbosity, capture request body (log later after handler)
		var reqBodyData []byte
		if s.verbose >= 2 {
			if r.Body != nil {
				data, _ := io.ReadAll(r.Body)
				// Restore the body for downstream handlers
				r.Body = io.NopCloser(bytes.NewReader(data))
				reqBodyData = data
			}
		}

		handler.ServeHTTP(wrapped, r)

		duration := time.Since(start)
		log.Printf("%s %s -> %d in %.1fms",
			r.Method,
			r.URL.Path,
			wrapped.statusCode,
			float64(duration.Microseconds())/1000.0,
		)

		if s.verbose >= 2 {
			if len(reqBodyData) > 0 {
				log.Printf("Request body:%s", formatMaybeJSON(reqBodyData))
			}

			respBody := wrapped.body.Bytes()
			if len(respBody) > 0 {
				log.Printf("Response body:%s", formatMaybeJSON(respBody))
			}
		}
	})
}

func formatMaybeJSON(data []byte) string {
	trimmed := bytes.TrimSpace(data)
	if len(trimmed) == 0 {
		return ""
	}
	if trimmed[0] == '{' || trimmed[0] == '[' {
		var v any
		if err := json.Unmarshal(trimmed, &v); err == nil {
			pretty, err := json.MarshalIndent(v, "", "  ")
			if err == nil {
				return "\n" + string(pretty)
			}
		}
	}
	return " " + string(data)
}
