package persistence

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/n0en0o/marketplace/internal/basket/domain"
	"github.com/n0en0o/marketplace/internal/shared"
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

func (r *CartRepository) Get(
	ctx context.Context,
	accountName string,
) (*domain.ShoppingCart, error) {
	var exists bool
	err := r.db.QueryRowContext(
		ctx,
		`SELECT EXISTS(
			SELECT 1 FROM shopping_carts
			WHERE account_name = $1
		)`,
		accountName,
	).Scan(&exists)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, shared.NewNotFoundError("Shopping Cart", accountName)
	}

	rows, err := r.db.QueryContext(
		ctx,
		`SELECT item_id, quantity, unit_price, item_title, item_note
		FROM shopping_cart_items
		WHERE account_name = $1`,
		accountName,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.ShoppingCartItem
	for rows.Next() {
		var item domain.ShoppingCartItem
		err := rows.Scan(
			&item.ItemID,
			&item.Quantity,
			&item.UnitPrice,
			&item.ItemTitle,
			&item.ItemNote,
		)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &domain.ShoppingCart{
		AccountName: accountName,
		Items:       items,
	}, nil
}

func (r *CartRepository) Remove(
	ctx context.Context,
	accountName string,
) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(
		ctx,
		`SELECT EXISTS(
			SELECT 1 FROM shopping_carts
			WHERE account_name = $1
		)`,
		accountName,
	).Scan(&exists)
	if err != nil {
		return false, err
	}

	if !exists {
		return false, shared.NewNotFoundError("ShoppingCart", accountName)
	}

	_, err = r.db.ExecContext(
		ctx,
		`DELETE FROM shopping_carts
		WHERE account_name = $1`,
		accountName,
	)
	if err != nil {
		return false, err
	}

	return true, nil
}
