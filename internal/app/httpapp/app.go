package httpapp

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"strconv"
	"time"
)

type App struct {
	logger *slog.Logger
	Server *http.Server
}

func New(
	logger *slog.Logger,
	host string,
	port int,
	timeout time.Duration,
	idleTimeout time.Duration,
	router http.Handler,
) *App {
	server := &http.Server{
		Addr:         net.JoinHostPort(host, strconv.Itoa(port)),
		Handler:      router,
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
		IdleTimeout:  idleTimeout,
	}
	return &App{
		logger: logger,
		Server: server,
	}
}

func (a *App) MustStart() {
	const op = "httpapp.MustStart"
	a.logger.Info("server is listening", slog.String("addr", a.Server.Addr))
	if err := a.Server.ListenAndServe(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			a.logger.Info("server closed received", slog.String("op", op))
			return
		}
		panic(err)
	}
}

func (a *App) GracefulStop() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()
	a.Server.Shutdown(ctx)
}
