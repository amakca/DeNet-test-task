package app

import (
	"log/slog"
	"os"
	"time"
)

func configureLogging() {
	if os.Getenv("ENV") == "dev" || os.Getenv("ENV") == "development" {
		handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug, ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.String(slog.TimeKey, time.Now().Format(time.RFC3339))
			}
			return a
		}})
		slog.SetDefault(slog.New(handler))
		return
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo, ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey {
			return slog.String(slog.TimeKey, time.Now().Format(time.RFC3339))
		}
		return a
	}})
	slog.SetDefault(slog.New(handler))
}
