package commands

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/n0en0o/marketplace/internal/catalog/domain/entities"
	"github.com/n0en0o/marketplace/internal/catalog/domain/repositories"
)

type UpdateCatalogItemCommand struct {
	ID               uuid.UUID          `json:"id"`
	Title            string             `json:"title"`
	ShortDescription string             `json:"short_description"`
	FullDescription  string             `json:"full_description"`
	ImageURL         string             `json:"image_url"`
	Price            float64            `json:"price"`
	BrandID          *entities.Brand    `json:"brand"`
	CategoryID       *entities.Category `json:"category"`
}

type UpdateCatalogItemHandler struct {
	repo repositories.CatalogItemRepository
}

func NewUpdateCatalogItemHandler(
	repo repositories.CatalogItemRepository,
) *UpdateCatalogItemHandler {
	return &UpdateCatalogItemHandler{
		repo: repo,
	}
}

func (h *UpdateCatalogItemHandler) Handle(
	ctx context.Context,
	cmd UpdateCatalogItemCommand,
) (bool, error) {

	if _, err := h.repo.Item(ctx, cmd.ID); errors.Is(err, repositories.ErrItemNotFound) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	item := entities.CatalogItem{
		BaseEntity: entities.BaseEntity{
			ID:    cmd.ID,
			Title: cmd.Title,
		},
		ShortDescription: cmd.ShortDescription,
		FullDescription:  cmd.FullDescription,
		ImageURL:         cmd.ImageURL,
		Price:            cmd.Price,
		Brand:            cmd.BrandID,
		Category:         cmd.CategoryID,
	}

	return h.repo.Update(ctx, item)
}
