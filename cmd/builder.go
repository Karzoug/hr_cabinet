package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/Employee-s-file-cabinet/backend/internal/config"
	"github.com/Employee-s-file-cabinet/backend/pkg/logger/slog/pretty"
)

func buildLogger(level string, envMode config.EnvType) *slog.Logger {
	var (
		logger    *slog.Logger
		slogLevel slog.Level
	)

	switch level {
	case "info":
		slogLevel = slog.LevelInfo
	case "error":
		slogLevel = slog.LevelError
	case "debug":
	default:
		slogLevel = slog.LevelDebug
	}

	switch envMode { // nolint:exhaustive
	case config.EnvDevelopment:
		opts := pretty.HandlerOptions{
			SlogOpts: &slog.HandlerOptions{
				Level: slogLevel,
			},
		}

		handler := opts.NewPrettyHandler(os.Stdout)
		logger = slog.New(handler)
	default:
		logger = slog.New(
			slog.NewJSONHandler(os.Stdout,
				&slog.HandlerOptions{
					Level:       slogLevel,
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
