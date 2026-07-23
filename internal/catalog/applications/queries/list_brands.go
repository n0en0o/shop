package queries

import (
	"context"

	"github.com/n0en0o/marketplace/internal/catalog/domain/entities"
	"github.com/n0en0o/marketplace/internal/catalog/domain/repositories"
)

type BrandsHandler struct {
	repo repositories.BrandRepository
}

func NewBrandsHandler(repo repositories.BrandRepository) *BrandsHandler {
	return &BrandsHandler{
		repo: repo,
	}
}

func (h *BrandsHandler) Handle(ctx context.Context) ([]entities.Brand, error) {
	return h.repo.Brands(ctx)
}
