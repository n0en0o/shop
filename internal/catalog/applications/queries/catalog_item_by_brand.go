package queries

import (
	"context"

	"github.com/n0en0o/shop/internal/catalog/domain/entities"
	"github.com/n0en0o/shop/internal/catalog/domain/repositories"
)

type CatalogItemsByBrandHandler struct {
	repo repositories.CatalogItemRepository
}

func NewCatalogItemsByBrandHandler(
	repo repositories.CatalogItemRepository,
) *CatalogItemsByBrandHandler {
	return &CatalogItemsByBrandHandler{
		repo: repo,
	}
}

func (h *CatalogItemsByBrandHandler) Handle(
	ctx context.Context,
	brand string,
) ([]entities.CatalogItem, error) {
	return h.repo.ItemsByBrand(ctx, brand)
}
