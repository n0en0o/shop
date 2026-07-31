package queries

import (
	"context"

	"github.com/n0en0o/marketplace/internal/promotion/applications/interfaces"
	"github.com/n0en0o/marketplace/internal/promotion/domain"
)

type GetByCatalogItemQuery struct {
	CatalogItemID string
}

type GetByCatalogItemHandler struct {
	repo interfaces.PromotionRepository
}

func NewGetByCatalogItemHandler(
	repo interfaces.PromotionRepository,
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
