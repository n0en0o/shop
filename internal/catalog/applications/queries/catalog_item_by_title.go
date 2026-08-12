package queries

import (
	"context"

	"github.com/n0en0o/shop/internal/catalog/domain/entities"
	"github.com/n0en0o/shop/internal/catalog/domain/repositories"
)

type CatalogItemsByTitleHandler struct {
	repo repositories.CatalogItemRepository
}

func NewCatalogItemsByTitleHandler(
	repo repositories.CatalogItemRepository) *CatalogItemsByTitleHandler {
	return &CatalogItemsByTitleHandler{
		repo: repo,
	}
}

func (h *CatalogItemsByTitleHandler) Handle(
	ctx context.Context, title string) ([]entities.CatalogItem, error) {
	return h.repo.ItemsByTitle(ctx, title)
}
