package commands

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/n0en0o/shop/internal/catalog/domain/repositories"
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
	if _, err := h.repo.Item(ctx, cmd.ID); errors.Is(err, repositories.ErrItemNotFound) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return h.repo.Delete(ctx, cmd.ID)
}
