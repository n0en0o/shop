package interfaces

import (
	"context"

	"github.com/google/uuid"
	"github.com/n0en0o/shop/internal/checkout/domain"
)

type OrderRepository interface {
	Get(ctx context.Context, id uuid.UUID) (*domain.Order, error)
	FindByAccountName(ctx context.Context, accountName string) ([]*domain.Order, error)
	Create(ctx context.Context, order *domain.Order) error
}
