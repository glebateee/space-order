package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/glebateee/space-order/internal/app"
	"github.com/glebateee/space-order/internal/config"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()
	logger := setupLogger(cfg.Env)
	logger.Info("loaded configuration", slog.Any("cfg", cfg))

	mainApp := app.New(
		logger,
		cfg.HttpConfig.Host,
		cfg.HttpConfig.Port,
		cfg.HttpConfig.Timeout,
		cfg.HttpConfig.IdleTimeout,
	)

	go mainApp.HttpApp.MustStart()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	sign := <-stop
	logger.Info("shutting down application", slog.String("signal", sign.String()))
	mainApp.HttpApp.GracefulStop()
	logger.Info("application stopped")

}

func setupLogger(env string) *slog.Logger {
	var logger *slog.Logger
	switch env {
	case envLocal:
		logger = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		logger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		logger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return logger
}
