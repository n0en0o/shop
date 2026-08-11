package repositories

import (
	"context"
	"errors"

	"github.com/n0en0o/marketplace/internal/promotion/domain"
)

var ErrPromoNotFound = errors.New("promotion not found")

type PromotionRepository interface {
	FindByCatalogItem(ctx context.Context, catalogItemID string) (*domain.Promo, error)
	Create(ctx context.Context, promo *domain.Promo) (bool, error)
	Update(ctx context.Context, promo *domain.Promo) (bool, error)
	Delete(ctx context.Context, id string) (bool, error)
}
