package middleware

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"lfg/shared/logging"
	"lfg/shared/metrics"
)

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	written    int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	rw.written += n
	return n, err
}

// RequestLogging creates request logging middleware
func RequestLogging(logger *logging.Logger, serviceName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Generate request ID if not present
			requestID := r.Header.Get("X-Request-ID")
			if requestID == "" {
				requestID = uuid.New().String()
			}

			// Add request ID to response header
			w.Header().Set("X-Request-ID", requestID)

			// Create logger with request context
			reqLogger := logger.
				WithRequestID(requestID).
				WithField("method", r.Method).
				WithField("path", r.URL.Path).
				WithField("remote_addr", r.RemoteAddr)

			// Add request ID to context
			ctx := logging.ContextWithRequestID(r.Context(), requestID)
			r = r.WithContext(ctx)

			// Wrap response writer to capture status code
			rw := newResponseWriter(w)

			// Log request start
			reqLogger.Infof("Request started: %s %s", r.Method, r.URL.Path)

			// Process request
			next.ServeHTTP(rw, r)

			// Calculate duration
			duration := time.Since(start)

			// Log request completion
			reqLogger.
				WithField("status", rw.statusCode).
				WithField("duration_ms", duration.Milliseconds()).
				WithField("bytes", rw.written).
				Infof("Request completed: %s %s - %d (%dms)",
					r.Method,
					r.URL.Path,
					rw.statusCode,
					duration.Milliseconds(),
				)

			// Record metrics
			metrics.RecordHTTPRequest(
				serviceName,
				r.Method,
				r.URL.Path,
				http.StatusText(rw.statusCode),
				duration,
			)
		})
	}
}
