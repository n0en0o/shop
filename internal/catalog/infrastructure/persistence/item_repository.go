package persistence

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/n0en0o/marketplace/internal/catalog/domain/entities"
)

const sqlCatalogItemsQuery = `
		SELECT ci.id,
			ci.title,
			ci.short_description,
			ci.full_description,
			ci.image_url,
			ci.price,
			b.id,
			b.title,
			c.id,
			c.title
		FROM catalog_items ci
		LEFT JOIN brands b ON ci.brand_id = b.id
		LEFT JOIN categories c ON ci.category_id = c.id`

type itemRepository struct {
	db *sql.DB
}

func NewItemRepository(db *sql.DB) *itemRepository {
	return &itemRepository{
		db: db,
	}
}

func (r *itemRepository) Items(ctx context.Context) ([]entities.CatalogItem, error) {
	query := sqlCatalogItemsQuery

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []entities.CatalogItem

	for rows.Next() {
		item, err := ScanCatalogItem(rows)
		if err != nil {
			return nil, err
		}

		items = append(items, *item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

type scanner interface {
	Scan(dest ...any) error
}

func ScanCatalogItem(s scanner) (*entities.CatalogItem, error) {
	var item *entities.CatalogItem = &entities.CatalogItem{}
	var shortDescription *sql.NullString
	var fullDescription *sql.NullString
	var imageURL *sql.NullString
	var brandID *uuid.NullUUID
	var brandTitle *sql.NullString
	var categoryID *uuid.NullUUID
	var categoryTitle *sql.NullString

	if err := s.Scan(
		&item.ID,
		&item.Title,
		&shortDescription,
		&fullDescription,
		&imageURL,
		&item.Price,
		&brandID,
		&brandTitle,
		&categoryID,
		&categoryTitle,
	); err != nil {
		return &entities.CatalogItem{}, err
	}

	item.ShortDescription = shortDescription.String
	item.FullDescription = fullDescription.String
	item.ImageURL = imageURL.String

	if brandID.Valid {
		item.Brand = &entities.Brand{
			BaseEntity: entities.BaseEntity{
				ID:    brandID.UUID,
				Title: brandTitle.String,
			},
		}
	}

	if categoryID.Valid {
		item.Category = &entities.Category{
			BaseEntity: entities.BaseEntity{
				ID:    categoryID.UUID,
				Title: categoryTitle.String,
			},
		}
	}

	return item, nil
}
