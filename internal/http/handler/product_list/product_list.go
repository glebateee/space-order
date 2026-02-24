package productlist

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/glebateee/space-order/internal/http/handler/product"
	"github.com/glebateee/space-order/internal/lib/sl"
	"github.com/glebateee/space-order/internal/service"
)

func New(
	ctx context.Context,
	logger *slog.Logger,
	provider product.Provider,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		list, err := provider.ProductList(ctx)
		if err != nil {
			if errors.Is(err, service.ErrInvalid) {
				logger.Info("invalid request", sl.Err(err))
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			logger.Info("internal error", sl.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(list)
	}
}
