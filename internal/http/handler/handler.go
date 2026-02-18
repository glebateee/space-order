package handler

import (
	"log/slog"
	"net/http"

	mwlogger "github.com/glebateee/space-order/internal/http/middleware/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func New(logger *slog.Logger) http.Handler {
	router := chi.NewRouter()
	router.Use(
		middleware.Recoverer,
		mwlogger.New(logger),
	)
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello"))
	})
	return router
}
