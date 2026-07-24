package queries

import (
	"context"

	"github.com/google/uuid"
	"github.com/n0en0o/marketplace/internal/catalog/domain/entities"
	"github.com/n0en0o/marketplace/internal/catalog/domain/repositories"
)

type CatalogItemByIdHandler struct {
	repo repositories.CatalogItemRepository
}

func NewCatalogItemByIdHandler(repo repositories.CatalogItemRepository) *CatalogItemByIdHandler {
	return &CatalogItemByIdHandler{
		repo: repo,
	}
}

func (h *CatalogItemByIdHandler) Handle(ctx context.Context, id uuid.UUID) (*entities.CatalogItem, error) {
	return h.repo.Item(ctx, id)
}
