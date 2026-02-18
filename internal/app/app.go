package app

import (
	"log/slog"
	"time"

	"github.com/glebateee/space-order/internal/app/httpapp"
	"github.com/glebateee/space-order/internal/http/handler"
)

type App struct {
	HttpApp *httpapp.App
}

func New(logger *slog.Logger, host string, port int, timeout time.Duration, idleTimeut time.Duration) *App {
	router := handler.New(logger)
	httpApp := httpapp.New(logger, host, port, timeout, idleTimeut, router)
	return &App{
		HttpApp: httpApp,
	}
}
