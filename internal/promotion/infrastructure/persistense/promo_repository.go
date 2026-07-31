package persistense

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
