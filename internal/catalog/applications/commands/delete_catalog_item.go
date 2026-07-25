package commands

import (
	"context"

	"github.com/google/uuid"
	"github.com/n0en0o/marketplace/internal/catalog/domain/repositories"
)

type DeleteCatalogItemCommand struct {
	ID uuid.UUID `json:"id"`
}

type DeleteCatalogItemHandler struct {
	repo repositories.CatalogItemRepository
}

func NewDeleteCatalogItemHandler(
	repo repositories.CatalogItemRepository,
) *DeleteCatalogItemHandler {
	return &DeleteCatalogItemHandler{
		repo: repo,
	}
}

func (h *DeleteCatalogItemHandler) Handle(
	ctx context.Context,
	cmd DeleteCatalogItemCommand,
) (bool, error) {
	existingItem, err := h.repo.Item(ctx, cmd.ID)

	if err != nil {
		return false, err
	}

	if existingItem == nil {
		return false, nil
	}

	return h.repo.Delete(ctx, cmd.ID)
}
