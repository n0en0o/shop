package persistence

import (
	"context"
	"database/sql"

	"github.com/n0en0o/marketplace/internal/catalog/domain/entities"
	"github.com/n0en0o/marketplace/internal/catalog/domain/repositories"
)

type categoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) repositories.CategoryRepository {
	return &categoryRepository{
		db: db,
	}
}

func (r *categoryRepository) Categories(ctx context.Context) ([]entities.Category, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, title FROM categories ORDER BY title ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []entities.Category

	for rows.Next() {
		var cat entities.Category
		if err := rows.Scan(&cat.ID, &cat.Title); err != nil {
			return nil, err
		}
		categories = append(categories, cat)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}
