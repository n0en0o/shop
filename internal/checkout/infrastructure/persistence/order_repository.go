package persistence

import (
	"context"
	"database/sql"

	"github.com/google/uuid"

	"github.com/n0en0o/shop/internal/checkout/applications/interfaces"
	"github.com/n0en0o/shop/internal/checkout/domain"
	"github.com/n0en0o/shop/internal/shared"
)

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) interfaces.OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Get(
	ctx context.Context,
	id uuid.UUID,
) (*domain.Order, error) {
	query := `
	SELECT
		id,
		account_name,
		total_amount,
		current_order_status,
		contact_first_name,
		contact_last_name,
		contact_email,
		address_street,
		address_city,
		address_region,
		address_postal_code,
		current_payment_method,
		current_payment_status,
		card_name,
		card_number,
		card_expiration,
		card_cvv,
		created_by,
		created_at,
		last_modified_by,
		last_modified_at
	FROM orders
	WHERE id = $1
	`

	orders, err := r.queryOrders(ctx, query, id)
	if err != nil {
		return nil, err
	}
	if len(orders) == 0 {
		return nil, shared.NewNotFoundError("Order", id)
	}

	return orders[0], nil
}

func (r *OrderRepository) getOrderItems(
	ctx context.Context,
	orderID uuid.UUID,
) ([]domain.OrderItem, error) {
	const query = `
		SELECT catalog_item_name, quantity, unit_price
		FROM order_items
		WHERE order_id = $1
		`
	rows, err := r.db.QueryContext(ctx, query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.OrderItem
	for rows.Next() {
		var item domain.OrderItem
		if err := rows.Scan(
			&item.CatalogItemName,
			&item.Quantity,
			&item.UnitPrice,
		); err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	return items, rows.Err()
}

func (r *OrderRepository) FindByAccountName(ctx context.Context, accountName string) ([]*domain.Order, error) {
	query := `
	SELECT
		id,
		account_name,
		total_amount,
		current_order_status,
		contact_first_name,
		contact_last_name,
		contact_email,
		address_street,
		address_city,
		address_region,
		address_postal_code,
		current_payment_method,
		current_payment_status,
		card_name,
		card_number,
		card_expiration,
		card_cvv,
		created_by,
		created_at,
		last_modified_by,
		last_modified_at
	FROM orders
	WHERE account_name = $1
	ORDER BY created_at DESC
	`

	return r.queryOrders(ctx, query, accountName)
}

func (r *OrderRepository) queryOrders(
	ctx context.Context,
	query string,
	args ...any,
) ([]*domain.Order, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*domain.Order
	for rows.Next() {
		order := &domain.Order{}
		var cardName, cardNumber, cardExp, cardCVV sql.NullString

		err := rows.Scan(
			&order.ID,
			&order.AccountName,
			&order.TotalAmount,
			&order.CurrentOrderStatus,
			&order.ContactInfo.FirstName,
			&order.ContactInfo.LastName,
			&order.ContactInfo.Email,
			&order.DeliveryAddress.Street,
			&order.DeliveryAddress.City,
			&order.DeliveryAddress.Region,
			&order.DeliveryAddress.PostalCode,
			&order.CurrentPaymentMethod,
			&order.CurrentPaymentStatus,
			&cardName,
			&cardNumber,
			&cardExp,
			&cardCVV,
			&order.CreatedBy,
			&order.CreatedAt,
			&order.LastModifiedBy,
			&order.LastModifiedAt,
		)
		if err != nil {
			return nil, err
		}

		if cardName.Valid || cardNumber.Valid || cardExp.Valid || cardCVV.Valid {
			order.CardDetails = &domain.CardDetails{
				CardName:   cardName.String,
				CardNumber: cardNumber.String,
				Expiration: cardExp.String,
				CVV:        cardCVV.String,
			}
		}

		items, err := r.getOrderItems(ctx, order.ID)
		if err != nil {
			return nil, err
		}
		order.Items = items

		orders = append(orders, order)
	}

	return orders, rows.Err()
}

func (r *OrderRepository) Create(
	ctx context.Context,
	order *domain.Order,
) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	const orderQuery = `
	INSERT INTO orders (
		id,
		account_name,
		total_amount,
		current_order_status,
		contact_first_name,
		contact_last_name,
		contact_email,
		address_street,
		address_city,
		address_region,
		address_postal_code,
		current_payment_method,
		current_payment_status,
		card_name,
		card_number,
		card_expiration,
		card_cvv,
		created_by,
		created_at
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7,
		$8, $9, $10, $11, $12, $13,
		$14, $15, $16, $17, $18, $19
	)
	`

	var cardName, cardNumber, cardExp, cardCVV *string
	if order.CardDetails != nil {
		cardName = &order.CardDetails.CardName
		cardNumber = &order.CardDetails.CardNumber
		cardExp = &order.CardDetails.Expiration
		cardCVV = &order.CardDetails.CVV
	}

	_, err = tx.ExecContext(
		ctx,
		orderQuery,
		order.ID,
		order.AccountName,
		order.TotalAmount,
		order.CurrentOrderStatus,
		order.ContactInfo.FirstName,
		order.ContactInfo.LastName,
		order.ContactInfo.Email,
		order.DeliveryAddress.Street,
		order.DeliveryAddress.City,
		order.DeliveryAddress.Region,
		order.DeliveryAddress.PostalCode,
		order.CurrentPaymentMethod,
		order.CurrentPaymentStatus,
		cardName,
		cardNumber,
		cardExp,
		cardCVV,
		order.CreatedBy,
		order.CreatedAt,
	)
	if err != nil {
		return err
	}

	const itemQuery = `
	INSERT INTO order_items (
		order_id,
		catalog_item_name,
		quantity,
		unit_price
	) VALUES ($1, $2, $3, $4)
	`

	for _, item := range order.Items {
		_, err := tx.ExecContext(
			ctx,
			itemQuery,
			order.ID,
			item.CatalogItemName,
			item.Quantity,
			item.UnitPrice,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
