package persistence

import (
	"context"
	"database/sql"

	"github.com/n0en0o/marketplace/internal/promotion/applications/interfaces"
	"github.com/n0en0o/marketplace/internal/promotion/domain"
)

type PromoRepository struct {
	db *sql.DB
}

func NewPromoRepository(db *sql.DB) interfaces.PromotionRepository {
	return &PromoRepository{db: db}
}

func (r *PromoRepository) FindByCatalogItem(
	ctx context.Context,
	catalogItemID string,
) (*domain.Promo, error) {
	query := `
	SELECT
		id,
		catalog_item_id,
		title,
		value
	FROM promos
	WHERE catalog_item_id = ?
	LIMIT 1
	`

	var p domain.Promo
	err := r.db.QueryRowContext(ctx, query, catalogItemID).Scan(
		&p.ID,
		&p.CatalogItemID,
		&p.Title,
		&p.Value,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (r *PromoRepository) Create(
	ctx context.Context,
	promo *domain.Promo,
) (bool, error) {
	const query = `
	INSERT INTO promos (id, catalog_item_id, title, value)
	VALUES (?, ?, ?, ?)
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		promo.ID,
		promo.CatalogItemID,
		promo.Title,
		promo.Value,
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
