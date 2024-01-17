package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Employee-s-file-cabinet/backend/internal/app"
	"github.com/Employee-s-file-cabinet/backend/internal/config"
)

// nolint:gochecknoglobals
var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		slog.Error(fmt.Sprintf("failed to get config: %s", err.Error()))
		os.Exit(1)
	}

	logger := buildLogger(cfg.LogLevel, cfg.EnvType)

	logger.Info(
		"starting app",
		slog.String("environment type", cfg.EnvType.String()),
		slog.String("build version", buildVersion),
		slog.String("build date", buildDate),
		slog.String("build commit", buildCommit),
	)

	ctx, stop := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	if err := app.Run(ctx, cfg, logger); err != nil {
		logger.Error("app stopped with error", slog.Attr{
			Key:   "error",
			Value: slog.StringValue(err.Error()),
		})
	}
}
