package commands

import (
	"context"

	"github.com/google/uuid"
	"github.com/n0en0o/marketplace/internal/catalog/domain/entities"
	"github.com/n0en0o/marketplace/internal/catalog/domain/repositories"
)

type CreateCatalogItemCommand struct {
	Title            string             `json:"title"`
	ShortDescription string             `json:"short_description"`
	FullDescription  string             `json:"full_description"`
	ImageURL         string             `json:"image_url"`
	Price            float64            `json:"price"`
	BrandID          *entities.Brand    `json:"brand"`
	CategoryID       *entities.Category `json:"category"`
}

type CreateCatalogItemHandler struct {
	repo repositories.CatalogItemRepository
}

func NewCreateCatalogItemHandler(
	repo repositories.CatalogItemRepository,
) *CreateCatalogItemHandler {
	return &CreateCatalogItemHandler{
		repo: repo,
	}
}

func (h *CreateCatalogItemHandler) Handle(
	ctx context.Context,
	cmd CreateCatalogItemCommand,
) (uuid.UUID, error) {
	item := entities.CatalogItem{
		BaseEntity: entities.BaseEntity{
			ID:    uuid.New(),
			Title: cmd.Title,
		},
		ShortDescription: cmd.ShortDescription,
		FullDescription:  cmd.FullDescription,
		ImageURL:         cmd.ImageURL,
		Price:            cmd.Price,
		Brand:            cmd.BrandID,
		Category:         cmd.CategoryID,
	}

	createdItem, err := h.repo.Create(ctx, item)
	if err != nil {
		return uuid.Nil, err
	}

	return createdItem.ID, nil
}
