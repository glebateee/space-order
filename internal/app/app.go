package app

import (
	"context"
	"log/slog"
	"time"

	"github.com/glebateee/space-order/internal/app/httpapp"
	"github.com/glebateee/space-order/internal/http/handler"
	"github.com/glebateee/space-order/internal/service/inventory"
)

type App struct {
	HttpApp *httpapp.App
}

func New(
	logger *slog.Logger,
	host string,
	port int,
	timeout time.Duration,
	idleTimeut time.Duration,
	grpcHost string,
	gRpcPort int,
) *App {
	invClient, err := inventory.New(logger, grpcHost, gRpcPort)
	if err != nil {
		panic(err)
	}
	handler := handler.New(
		context.Background(),
		logger,
		invClient,
	)
	httpApp := httpapp.New(
		logger,
		host,
		port,
		timeout,
		idleTimeut,
		handler,
	)
	return &App{
		HttpApp: httpApp,
	}
}
