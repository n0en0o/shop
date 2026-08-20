package persistence

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/n0en0o/shop/internal/checkout/domain"

	"github.com/n0en0o/shop/internal/checkout/applications/interfaces"
)

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) interfaces.OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Get(ctx context.Context, id uuid.UUID) (*domain.Order, error) {
	return nil, nil
}

func (r *OrderRepository) FindByAccountName(ctx context.Context, accountName string) ([]*domain.Order, error) {
	return nil, nil
}

func (r *OrderRepository) Create(ctx context.Context, order *domain.Order) error {
	return nil
}
