package persistence

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/n0en0o/marketplace/internal/basket/domain"
)

type CartRepository struct {
	db *sql.DB
}

func NewCartRepository(db *sql.DB) *CartRepository {
	return &CartRepository{db: db}
}

func (r *CartRepository) Save(
	ctx context.Context,
	cart domain.ShoppingCart,
) (*domain.ShoppingCart, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(
		ctx,
		`INSERT INTO shopping_carts (account_name)
		VALUES ($1)
		ON CONFLICT (account_name) DO NOTHING`,
		cart.AccountName,
	)
	if err != nil {
		return nil, err
	}

	_, err = tx.ExecContext(
		ctx,
		`DELETE FROM shopping_cart_items
		WHERE account_name = $1`,
		cart.AccountName,
	)
	if err != nil {
		return nil, err
	}

	for _, item := range cart.Items {
		if item.ItemID == uuid.Nil {
			item.ItemID = uuid.New()
		}

		_, err = tx.ExecContext(
			ctx,
			`INSERT INTO shopping_cart_items (
				account_name,
				item_id,
				quantity,
				unit_price,
				item_title,
				item_note
			) VALUES ($1, $2, $3, $4, $5, $6)`,
			cart.AccountName,
			item.ItemID,
			item.Quantity,
			item.UnitPrice,
			item.ItemTitle,
			item.ItemNote,
		)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &cart, nil
}
