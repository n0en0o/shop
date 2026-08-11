package queries

import (
	"context"

	"github.com/n0en0o/marketplace/internal/promotion/domain"
	"github.com/n0en0o/marketplace/internal/promotion/domain/repositories"
)

type GetByCatalogItemQuery struct {
	CatalogItemID string
}

type GetByCatalogItemHandler struct {
	repo repositories.PromotionRepository
}

func NewGetByCatalogItemHandler(
	repo repositories.PromotionRepository,
) *GetByCatalogItemHandler {
	return &GetByCatalogItemHandler{
		repo: repo,
	}
}

func (h *GetByCatalogItemHandler) Handle(
	ctx context.Context,
	q GetByCatalogItemQuery,
) (*domain.Promo, error) {
	return h.repo.FindByCatalogItem(ctx, q.CatalogItemID)
}
