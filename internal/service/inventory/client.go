package inventory

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"strconv"

	"github.com/glebateee/space-order/internal/domain/models"
	"github.com/glebateee/space-order/internal/lib/sl"
	"github.com/glebateee/space-order/internal/service"
	invv1 "github.com/glebateee/space-proto/gen/go/inventory"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type Inventory struct {
	logger *slog.Logger
	api    invv1.InventoryClient
}

func New(
	logger *slog.Logger,
	host string,
	port int,
) (*Inventory, error) {
	const op = "inventory.New"
	cc, err := grpc.NewClient(
		net.JoinHostPort(host, strconv.Itoa(port)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("%s :%w", op, err)
	}
	api := invv1.NewInventoryClient(cc)
	return &Inventory{
		logger: logger,
		api:    api,
	}, nil
}

func (iv *Inventory) ProductSku(ctx context.Context, sku string) (*models.Product, error) {
	const op = "inventory.ProductSku"
	logger := iv.logger.With(
		slog.String("op", op),
	)
	resp, err := iv.api.GetProduct(ctx,
		&invv1.GetProductRequest{
			Sku: sku,
		})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			logger.Error("failed convert to grpc error", sl.Err(err))
			return nil, service.ErrInternal
		}
		switch st.Code() {
		case codes.InvalidArgument:
			logger.Error("failed to process", sl.Err(err))
			return nil, service.ErrInvalid
		default:
			logger.Error("failed to handle error type", sl.Err(err))
			return nil, service.ErrInternal
		}
	}
	respProd := resp.GetProduct()
	return gRPCToProduct(respProd), nil
}

func (iv *Inventory) ProductList(ctx context.Context) ([]*models.Product, error) {
	const op = "inventory.ProductList"
	logger := iv.logger.With(
		slog.String("op", op),
	)
	resp, err := iv.api.GetProductList(ctx, &invv1.GetProductListRequest{})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Info("no rows found")
			return nil, service.ErrInvalid
		}
		logger.Error("unexpected error", sl.Err(err))
		return nil, service.ErrInternal
	}
	respProds := resp.GetProducts()
	prods := make([]*models.Product, 0, len(respProds))
	for _, p := range respProds {
		prods = append(prods, gRPCToProduct(p))
	}
	return prods, nil
}
