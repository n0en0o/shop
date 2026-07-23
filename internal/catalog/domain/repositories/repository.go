package repositories

import (
	"context"

	"github.com/n0en0o/marketplace/internal/catalog/domain/entities"
)

type CatalogItemRepository interface {
	Items(ctx context.Context) ([]entities.CatalogItem, error)
}

type BrandRepository interface {
	Brands(ctx context.Context) ([]entities.Brand, error)
}

type CategoryRepository interface {
	Categories(ctx context.Context) ([]entities.Category, error)
}
