package persistence

import (
	"context"
	"database/sql"
	"github.com/n0en0o/marketplace/internal/catalog/domain/entities"
)

type brandRepository struct {
	db *sql.DB
}

func NewBrandRepository(db *sql.DB) *brandRepository {
	return &brandRepository{
		db: db,
	}
}

func (r *brandRepository) Brands(ctx context.Context) ([]entities.Brand, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, title FROM brands")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var brands []entities.Brand

	for rows.Next() {
		var brand entities.Brand
		if err := rows.Scan(&brand.ID, &brand.Title); err != nil {
			return nil, err
		}
		brands = append(brands, brand)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return brands, nil
}
