package main

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Employee-s-file-cabinet/backend/internal/app"
	"github.com/Employee-s-file-cabinet/backend/internal/config"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		slog.Error(fmt.Sprintf("failed to get config: %s", err.Error()))
		os.Exit(1)
	}

	slogInit(os.Stdout, cfg.LogLevel)

	a, err := app.New(cfg)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to init: %s", err.Error()))
		os.Exit(1)
	}

	err = a.Run()
	if err != nil {
		slog.Error(fmt.Sprintf("app run: %s", err.Error()))
	}

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	sig := <-interrupt
	slog.Info(fmt.Sprintf("server - Run - signal: %s", sig.String()))

	// TODO: Останов app или отмена контекста
}

func slogInit(out io.Writer, level string) {
	var slogLevel slog.Level

	switch level {
	case "info":
		slogLevel = slog.LevelInfo
	case "error":
		slogLevel = slog.LevelError
	case "debug":
	default:
		slogLevel = slog.LevelDebug
	}

	logger := slog.New(slog.NewJSONHandler(out,
		&slog.HandlerOptions{
			Level:       slogLevel,
			ReplaceAttr: rewriteSlogAttributes(),
		}),
	)
	slog.SetDefault(logger)
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
