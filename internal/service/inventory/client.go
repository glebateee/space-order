package inventory

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"strconv"

	"github.com/glebateee/space-order/internal/domain/models"
	invv1 "github.com/glebateee/space-proto/gen/go/inventory"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

func gRPCToProduct(p *invv1.Product) *models.Product {
	return &models.Product{
		UUID:        p.GetUuid(),
		SKU:         p.GetSku(),
		Name:        p.GetName(),
		Description: p.GetDescription(),
		Category:    p.GetCategory(),
		Currency:    p.GetCurrency(),
		BasePrice:   p.GetBasePrice(),
		CreatedAt:   p.GetCreatedAt().AsTime(),
		UpdatedAt:   p.GetUpdatedAt().AsTime(),
	}
}
func (iv *Inventory) ProductSku(ctx context.Context, sku string) (*models.Product, error) {
	resp, err := iv.api.GetProduct(ctx,
		&invv1.GetProductRequest{
			Sku: sku,
		})
	if err != nil {
		return nil, err
	}
	respProd := resp.GetProduct()
	return gRPCToProduct(respProd), nil
}
