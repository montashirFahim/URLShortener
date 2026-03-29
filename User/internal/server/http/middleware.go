package http

import (
	"log"
	"net/http"
	"time"
)

// Middleware represents http handler middleware
type Middleware func(http.Handler) http.Handler

func mergeMiddleware(h http.Handler, middlewares ...Middleware) http.Handler {
	for _, middleware := range middlewares {
		h = middleware(h)
	}

	return h
}

func closeBody(h http.Handler) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(handler)
}

// DetailedLogger is a middleware that logs detailed information about incoming requests
func DetailedLogger(h http.Handler) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Log detailed request information
		log.Printf("[REQUEST] %s %s %s %s %s",
			r.Method,
			r.URL.Path,
			r.URL.RawQuery,
			r.RemoteAddr,
			r.UserAgent(),
		)

		// Wrap the ResponseWriter to capture status code
		wrappedWriter := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Call the next handler
		h.ServeHTTP(wrappedWriter, r)

		// Log response information
		duration := time.Since(start)
		log.Printf("[RESPONSE] %s %s %d %v",
			r.Method,
			r.URL.Path,
			wrappedWriter.statusCode,
			duration,
		)
	}

	return http.HandlerFunc(handler)
}

// responseWriter is a wrapper to capture the status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
