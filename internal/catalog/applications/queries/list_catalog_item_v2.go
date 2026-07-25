package queries

import (
	"context"

	"github.com/n0en0o/marketplace/internal/catalog/domain/entities"
	"github.com/n0en0o/marketplace/internal/catalog/domain/repositories"
	"github.com/n0en0o/marketplace/internal/catalog/domain/spec"
)

type CatalogItemsV2Handler struct {
	repo repositories.CatalogItemRepository
}

func NewCatalogItemsV2Handler(repo repositories.CatalogItemRepository) *CatalogItemsV2Handler {
	return &CatalogItemsV2Handler{
		repo: repo,
	}
}

func (h *CatalogItemsV2Handler) Handle(
	ctx context.Context,
	args spec.QueryArgs,
) (spec.Pagination[entities.CatalogItem], error) {
	return h.repo.CatalogItems(ctx, args)
}
