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

func (iv *Inventory) Get(ctx context.Context) ([]models.Product, error) {
	resp, err := iv.api.GetProductList(ctx,
		&invv1.GetProductListRequest{})
	if err != nil {
		return nil, err
	}
	respProds := resp.GetProducts()
	prods := make([]models.Product, 0, len(respProds))
	for _, p := range respProds {
		prods = append(prods, models.Product{
			UUID: p.GetUuid(),
			SKU:  p.GetSku(),
		})
	}
	return prods, nil
}
