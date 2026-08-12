package persistence

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/n0en0o/shop/internal/promotion/domain"
	"github.com/n0en0o/shop/internal/promotion/domain/repositories"
)

const mysqlDuplicateEntryCode uint16 = 1062

type PromoRepository struct {
	db *sql.DB
}

func NewPromoRepository(db *sql.DB) repositories.PromotionRepository {
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
	if errors.Is(err, sql.ErrNoRows) {
		return nil, repositories.ErrPromoNotFound
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
		var myErr *mysql.MySQLError
		if errors.As(err, &myErr) && myErr.Number == mysqlDuplicateEntryCode {
			return false, fmt.Errorf("%w: catalog_item_id %s",
				domain.ErrPromoAlreadyExists,
				promo.CatalogItemID,
			)
		}
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected > 0, nil
}

func (r *PromoRepository) Update(
	ctx context.Context,
	promo *domain.Promo,
) (bool, error) {
	const query = `
	UPDATE promos SET title = ?, value = ?
	WHERE id = ?
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		promo.Title,
		promo.Value,
		promo.ID,
	)
	if err != nil {
		return false, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	if rows > 0 {
		return true, nil
	}

	return r.exists(ctx, promo.ID)
}

func (r *PromoRepository) Delete(
	ctx context.Context,
	id string,
) (bool, error) {
	const query = `
	DELETE FROM promos
	WHERE id = ?
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		id,
	)
	if err != nil {
		return false, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rows > 0, nil
}

func (r *PromoRepository) exists(ctx context.Context, id string) (bool, error) {
	var exists int
	err := r.db.QueryRowContext(
		ctx,
		`SELECT EXISTS(
			SELECT 1 FROM promos
			WHERE id = ?
		)`,
		id,
	).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists > 0, nil
}
