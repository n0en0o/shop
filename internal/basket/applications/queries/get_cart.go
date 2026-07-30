package queries

import (
	"context"

	"github.com/n0en0o/marketplace/internal/basket/domain"
	"github.com/n0en0o/marketplace/internal/basket/domain/repositories"
)

type GetCartHandler struct {
	repo repositories.CartRepository
}

func NewGetCartHandler(
	repo repositories.CartRepository,
) *GetCartHandler {
	return &GetCartHandler{repo: repo}
}

func (h *GetCartHandler) Handle(
	ctx context.Context,
	accountName string,
) (*domain.ShoppingCart, error) {
	return h.repo.Get(ctx, accountName)
}
