package interfaces

import (
	"context"

	"github.com/n0en0o/marketplace/internal/basket/domain"
)

type CartRepository interface {
	Save(ctx context.Context, cart domain.ShoppingCart) (*domain.ShoppingCart, error)
	Get(ctx context.Context, accountName string) (*domain.ShoppingCart, error)
	Remove(ctx context.Context, accountName string) (bool, error)
}
