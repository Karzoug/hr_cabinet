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
	"github.com/Employee-s-file-cabinet/backend/pkg/logger/slog/sl"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"

	envMode = config.EnvProduction
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	cfg, err := config.New()
	if err != nil {
		slog.Error(fmt.Sprintf("failed to get config: %s", err.Error()))
		os.Exit(1)
	}

	logger := buildLogger(os.Stdout, cfg.LogLevel, envMode)

	logger.Info(
		"starting app server",
		slog.String("env", envMode.String()),
		slog.String("build version", buildVersion),
		slog.String("build date", buildDate),
	)
	if err := app.Run(ctx, cfg, logger); err != nil {
		logger.Error("app stopped with error", sl.Error(err))
	}
}
