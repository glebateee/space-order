package product

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/glebateee/space-order/internal/domain/models"
	"github.com/glebateee/space-order/internal/lib/sl"
	"github.com/glebateee/space-order/internal/service"
	"github.com/go-chi/chi/v5"
)

var emptySku = ""

type Provider interface {
	ProductSku(ctx context.Context, sku string) (*models.Product, error)
	ProductList(ctx context.Context) ([]*models.Product, error)
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
			switch {
			case errors.Is(err, service.ErrInternal):
				w.WriteHeader(http.StatusInternalServerError)
				return
			case errors.Is(err, service.ErrInvalid):
				w.WriteHeader(http.StatusBadRequest)
				return
			default:
				logger.Error("error receiving product", sl.Err(err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			// logger.Error("error receiving product", sl.Err(err))
			// w.WriteHeader(http.StatusInternalServerError)
			// json.NewEncoder(w).Encode(struct {
			// 	Status string `json:"status"`
			// 	Err    string `json:"error"`
			// }{
			// 	Status: "error",
			// 	Err:    err.Error(),
			// })
			// return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(prod)
		w.WriteHeader(http.StatusOK)
	}
}
