package middleware

import (
	"net/http"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// responseWriter captures the HTTP status code
type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

// GCPLoggingInterceptor maps standard HTTP requests automatically to JSON structures
// inherently compatible with Google Cloud Logging definitions (Severity level mappings).
func GCPLoggingInterceptor(next http.Handler) http.Handler {
	// Zerolog maps "level" to GCP's expected severity. We set standard configs first.
	zerolog.LevelFieldName = "severity"
	zerolog.TimestampFieldName = "timestamp"
	zerolog.MessageFieldName = "message"

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}
		
		// Map down to next handler
		next.ServeHTTP(rw, r)
		
		duration := time.Since(start)

		logger := log.Info()
		if rw.status >= 500 {
			logger = log.Error()
		} else if rw.status >= 400 {
			logger = log.Warn()
		}

		// Inject deterministic structured spans mirroring GCP payload architecture
		logger.
			Str("httpRequest.requestMethod", r.Method).
			Str("httpRequest.requestUrl", r.URL.String()).
			Str("httpRequest.userAgent", r.UserAgent()).
			Int("httpRequest.status", rw.status).
			Str("httpRequest.latency", duration.String()).
			Msg("HTTP request processed")
	})
}
