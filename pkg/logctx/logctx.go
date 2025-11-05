package logctx

import (
	"context"
	"log/slog"
)

type ctxLoggerKey struct{}

// WithLogger returns a new context holding the provided slog.Logger.
func WithLogger(ctx context.Context, l *slog.Logger) context.Context {
	return context.WithValue(ctx, ctxLoggerKey{}, l)
}

// FromContext returns slog.Logger from context or slog.Default() if absent.
func FromContext(ctx context.Context) *slog.Logger {
	if l, ok := ctx.Value(ctxLoggerKey{}).(*slog.Logger); ok && l != nil {
		return l
	}
	return slog.Default()
}
