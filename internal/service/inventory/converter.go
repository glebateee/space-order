package inventory

import (
	"github.com/glebateee/space-order/internal/domain/models"
	invv1 "github.com/glebateee/space-proto/gen/go/inventory"
)

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
