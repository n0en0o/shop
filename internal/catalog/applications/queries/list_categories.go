package queries

import (
	"context"

	"github.com/n0en0o/shop/internal/catalog/domain/entities"
	"github.com/n0en0o/shop/internal/catalog/domain/repositories"
)

type CategoriesHandler struct {
	repo repositories.CategoryRepository
}

func NewCategoriesHandler(repo repositories.CategoryRepository) *CategoriesHandler {
	return &CategoriesHandler{
		repo: repo,
	}
}

func (h *CategoriesHandler) Handle(ctx context.Context) ([]entities.Category, error) {
	return h.repo.Categories(ctx)
}
