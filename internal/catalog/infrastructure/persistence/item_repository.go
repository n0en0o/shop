package persistence

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/n0en0o/marketplace/internal/catalog/domain/entities"
	"github.com/n0en0o/marketplace/internal/catalog/domain/repositories"
	"github.com/n0en0o/marketplace/internal/catalog/domain/spec"
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

func NewItemRepository(db *sql.DB) repositories.CatalogItemRepository {
	return &itemRepository{
		db: db,
	}
}

func (r *itemRepository) Items(
	ctx context.Context,
) ([]entities.CatalogItem, error) {
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

func (r *itemRepository) Item(
	ctx context.Context, id uuid.UUID,
) (*entities.CatalogItem, error) {
	query := sqlCatalogItemsQuery + " WHERE ci.id = $1"

	row := r.db.QueryRowContext(ctx, query, id)

	item, err := ScanCatalogItem(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repositories.ErrItemNotFound
		}
		return nil, err
	}

	return item, nil
}

func (r *itemRepository) ItemsByTitle(
	ctx context.Context, title string,
) ([]entities.CatalogItem, error) {
	query := sqlCatalogItemsQuery + " WHERE ci.title ILIKE '%' || $1 || '%'"

	rows, err := r.db.QueryContext(ctx, query, title)

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

func (r *itemRepository) ItemsByBrand(
	ctx context.Context, brand string,
) ([]entities.CatalogItem, error) {
	query := sqlCatalogItemsQuery + " WHERE b.title ILIKE '%' || $1 || '%'"

	rows, err := r.db.QueryContext(ctx, query, brand)

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

func (r *itemRepository) CatalogItems(ctx context.Context, args spec.QueryArgs) (
	spec.Pagination[entities.CatalogItem], error) {

	args.Normalize()

	sqlBaseFrom := `
		FROM catalog_items ci
		LEFT JOIN brands b ON ci.brand_id = b.id
		LEFT JOIN categories c ON ci.category_id = c.id
	`

	var conditions []string
	var params []any
	paramIdx := 1

	brandID, err := args.ParseBrandID()
	if err != nil {
		return spec.Pagination[entities.CatalogItem]{},
			fmt.Errorf("invalid brand_id: %w", err)
	}

	categoryID, err := args.ParseCategoryID()
	if err != nil {
		return spec.Pagination[entities.CatalogItem]{},
			fmt.Errorf("invalid category_id: %w", err)
	}

	if brandID != nil {
		conditions = append(conditions, fmt.Sprintf("ci.brand_id = $%d", paramIdx))
		params = append(params, *brandID)
		paramIdx += 1
	}

	if categoryID != nil {
		conditions = append(conditions, fmt.Sprintf("ci.category_id = $%d", paramIdx))
		params = append(params, *categoryID)
		paramIdx += 1
	}

	if args.Search != nil && *args.Search != "" {
		conditions = append(conditions,
			fmt.Sprintf("ci.title ILIKE '%%' || $%d || '%%'", paramIdx))
		params = append(params, *args.Search)
		paramIdx += 1
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = " WHERE " + strings.Join(conditions, " AND ")
	}

	countQuery := "SELECT count(*)" + sqlBaseFrom + whereClause

	var totalCount int
	err = r.db.QueryRowContext(ctx, countQuery, params...).Scan(&totalCount)
	if err != nil {
		return spec.Pagination[entities.CatalogItem]{}, err
	}

	orderClause := ""
	if args.Sort != nil {
		switch strings.ToLower(*args.Sort) {
		case "price_asc":
			orderClause = " ORDER BY ci.price ASC"
		case "price_desc":
			orderClause = " ORDER BY ci.price DESC"
		case "title_asc":
			orderClause = " ORDER BY ci.title ASC"
		case "title_desc":
			orderClause = " ORDER BY ci.title DESC"
		}
	}

	offset := (args.PageIndex - 1) * args.PageSize
	paginationClause := fmt.Sprintf(" LIMIT $%d OFFSET $%d", paramIdx, paramIdx+1)
	paginationParams := append(params, args.PageSize, offset)

	sqlSelectFields := `
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
	`

	sqlFullQuery := sqlSelectFields + sqlBaseFrom + whereClause + orderClause + paginationClause

	rows, err := r.db.QueryContext(ctx, sqlFullQuery, paginationParams...)
	if err != nil {
		return spec.Pagination[entities.CatalogItem]{}, err
	}
	defer rows.Close()

	var items []entities.CatalogItem = []entities.CatalogItem{}

	for rows.Next() {
		item, err := ScanCatalogItem(rows)
		if err != nil {
			return spec.Pagination[entities.CatalogItem]{}, err
		}
		items = append(items, *item)
	}

	if err := rows.Err(); err != nil {
		return spec.Pagination[entities.CatalogItem]{}, err
	}

	return spec.Pagination[entities.CatalogItem]{
		PageIndex:  args.PageIndex,
		PageSize:   args.PageSize,
		TotalCount: totalCount,
		Items:      items,
	}, nil
}

func (r *itemRepository) Create(
	ctx context.Context, item entities.CatalogItem,
) (entities.CatalogItem, error) {

	if item.ID == uuid.Nil {
		item.ID = uuid.New()
	}

	var brandID, categoryID *uuid.UUID

	if item.Brand != nil {
		brandID = &item.Brand.ID
	}

	if item.Category != nil {
		categoryID = &item.Category.ID
	}

	sqlCreateQuery := `
		INSERT INTO catalog_items (
			id, 
			title, 
			short_description, 
			full_description, 
			image_url,
			price, 
			brand_id, 
			category_id
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.ExecContext(
		ctx,
		sqlCreateQuery,
		item.ID,
		item.Title,
		item.ShortDescription,
		item.FullDescription,
		item.ImageURL,
		item.Price,
		brandID,
		categoryID,
	)

	if err != nil {
		return entities.CatalogItem{}, err
	}

	return item, nil

}

func (r *itemRepository) Update(
	ctx context.Context, item entities.CatalogItem,
) (bool, error) {

	var brandID, categoryID *uuid.UUID

	if item.Brand != nil {
		brandID = &item.Brand.ID
	}

	if item.Category != nil {
		categoryID = &item.Category.ID
	}

	sqlUpdateQuery := `
		UPDATE catalog_items
		SET title = $1,
			short_description = $2,
			full_description = $3,
			image_url = $4,
			price = $5,
			brand_id = $6,
			category_id = $7
		WHERE id = $8
	`

	result, err := r.db.ExecContext(
		ctx,
		sqlUpdateQuery,
		item.Title,
		item.ShortDescription,
		item.FullDescription,
		item.ImageURL,
		item.Price,
		brandID,
		categoryID,
		item.ID,
	)

	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return false, err
	}

	return rowsAffected > 0, nil
}

func (r *itemRepository) Delete(
	ctx context.Context, id uuid.UUID,
) (bool, error) {

	sqlDeleteQuery := `
		DELETE FROM catalog_items
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, sqlDeleteQuery, id)

	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return false, err
	}

	return rowsAffected > 0, nil
}

type scanner interface {
	Scan(dest ...any) error
}

func ScanCatalogItem(s scanner) (*entities.CatalogItem, error) {
	var item *entities.CatalogItem = &entities.CatalogItem{}
	var shortDescription sql.NullString
	var fullDescription sql.NullString
	var imageURL sql.NullString
	var brandID uuid.NullUUID
	var brandTitle sql.NullString
	var categoryID uuid.NullUUID
	var categoryTitle sql.NullString

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
