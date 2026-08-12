package queries

import (
	"context"

	"github.com/google/uuid"
	"github.com/n0en0o/shop/internal/catalog/domain/entities"
	"github.com/n0en0o/shop/internal/catalog/domain/repositories"
)

type CatalogItemByIDHandler struct {
	repo repositories.CatalogItemRepository
}

func NewCatalogItemByIDHandler(repo repositories.CatalogItemRepository) *CatalogItemByIDHandler {
	return &CatalogItemByIDHandler{
		repo: repo,
	}
}

func (h *CatalogItemByIDHandler) Handle(ctx context.Context, id uuid.UUID) (*entities.CatalogItem, error) {
	return h.repo.Item(ctx, id)
}
