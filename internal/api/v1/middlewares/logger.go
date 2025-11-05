package middlewares

import (
	"log/slog"
	"net"
	"net/http"
	"strings"
	"time"

	"denet-test-task/pkg/logctx"

	chimw "github.com/go-chi/chi/v5/middleware"
)

// SlogRequestContext attaches a request-scoped slog.Logger to context.
func SlogRequestContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := chimw.GetReqID(r.Context())
		logger := slog.With(
			"req_id", reqID,
			"method", r.Method,
			"path", r.URL.Path,
			"remote_ip", getRealIP(r),
		)
		next.ServeHTTP(w, r.WithContext(logctx.WithLogger(r.Context(), logger)))
	})
}

// SlogAccessLogger writes a structured access log entry per request.
func SlogAccessLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := chimw.NewWrapResponseWriter(w, r.ProtoMajor)

		next.ServeHTTP(ww, r)

		logctx.FromContext(r.Context()).Info("http_request",
			"status", ww.Status(),
			"bytes", ww.BytesWritten(),
			"duration_ms", time.Since(start).Milliseconds(),
		)
	})
}

// getRealIP extracts client IP from common headers or RemoteAddr.
func getRealIP(r *http.Request) string {
	// X-Forwarded-For can be a comma-separated list; take the first non-empty
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		if len(parts) > 0 {
			ip := strings.TrimSpace(parts[0])
			if ip != "" {
				return ip
			}
		}
	}
	if xrip := r.Header.Get("X-Real-IP"); xrip != "" {
		return strings.TrimSpace(xrip)
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil && host != "" {
		return host
	}
	return r.RemoteAddr
}
