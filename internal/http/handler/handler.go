package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/glebateee/space-order/internal/domain/models"
	mwlogger "github.com/glebateee/space-order/internal/http/middleware/logger"
	"github.com/glebateee/space-order/internal/lib/sl"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type productProvider interface {
	Get(ctx context.Context) ([]models.Product, error)
}
type Handler struct {
	logger          *slog.Logger
	productProvider productProvider
	router          http.Handler
}

func New(
	logger *slog.Logger,
	prodProvider productProvider,
) *Handler {
	router := chi.NewRouter()
	router.Use(
		middleware.Recoverer,
		mwlogger.New(logger),
	)
	router.Get("/products", getProducts(logger, prodProvider))
	return &Handler{
		logger:          logger,
		productProvider: prodProvider,
		router:          router,
	}
}

func getProducts(
	logger *slog.Logger,
	prodProvider productProvider,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		prod, err := prodProvider.Get(ctx)
		if err != nil {
			logger.Error("error receiving product", sl.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(prod)
		w.WriteHeader(http.StatusOK)
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}
