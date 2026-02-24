package handler

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/glebateee/space-order/internal/http/handler/product"
	productlist "github.com/glebateee/space-order/internal/http/handler/product_list"
	mwlogger "github.com/glebateee/space-order/internal/http/middleware/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Handler struct {
	logger          *slog.Logger
	productProvider product.Provider
	router          http.Handler
}

func New(
	ctx context.Context,
	logger *slog.Logger,
	prodProvider product.Provider,
) *Handler {
	router := chi.NewRouter()
	router.Use(
		middleware.Recoverer,
		mwlogger.New(logger),
	)
	router.Get("/product/{sku}", product.New(ctx, logger, prodProvider))
	router.Get("/products", productlist.New(ctx, logger, prodProvider))

	return &Handler{
		logger:          logger,
		productProvider: prodProvider,
		router:          router,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}
