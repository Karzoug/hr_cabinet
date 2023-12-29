package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/Employee-s-file-cabinet/backend/internal/config/env"
	"github.com/Employee-s-file-cabinet/backend/pkg/logger/slog/pretty"
)

func buildLogger(level slog.Level, envMode env.Type) *slog.Logger {
	var logger *slog.Logger

	switch envMode { // nolint:exhaustive
	case env.Development:
		opts := pretty.HandlerOptions{
			SlogOpts: &slog.HandlerOptions{
				Level: level,
			},
		}

		handler := opts.NewPrettyHandler(os.Stdout)
		logger = slog.New(handler)
	default:
		logger = slog.New(
			slog.NewJSONHandler(os.Stdout,
				&slog.HandlerOptions{
					Level:       level,
					ReplaceAttr: rewriteSlogAttributes(),
				}),
		)
	}

	slog.SetDefault(logger)

	return logger
}

func rewriteSlogAttributes() func(groups []string, a slog.Attr) slog.Attr {
	return func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey {
			a.Key = "timestamp"
			a.Value = slog.StringValue(time.Now().Format("2006-01-02 15:04:05 UTC-07"))
		}
		return a
	}
}
