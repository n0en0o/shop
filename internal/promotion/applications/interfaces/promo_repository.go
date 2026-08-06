package interfaces

import (
	"context"

	"github.com/n0en0o/marketplace/internal/promotion/domain"
)

type PromotionRepository interface {
	FindByCatalogItem(ctx context.Context, catalogItemID string) (*domain.Promo, error)
	Create(ctx context.Context, promo *domain.Promo) (bool, error)
}
