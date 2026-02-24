package product

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/glebateee/space-order/internal/domain/models"
	"github.com/glebateee/space-order/internal/lib/sl"
	"github.com/go-chi/chi/v5"
)

var emptySku = ""

type Provider interface {
	ProductSku(ctx context.Context, sku string) (*models.Product, error)
}

func New(ctx context.Context, logger *slog.Logger, prodProvider Provider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sku := chi.URLParam(r, "sku")
		if sku == emptySku {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		prod, err := prodProvider.ProductSku(ctx, sku)
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
