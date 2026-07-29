package queries

import (
	"context"

	"github.com/n0en0o/marketplace/internal/basket/applications/interfaces"
	"github.com/n0en0o/marketplace/internal/basket/domain"
)

type GetCartHandler struct {
	repo interfaces.CartRepository
}

func NewGetCartHandler(
	repo interfaces.CartRepository,
) *GetCartHandler {
	return &GetCartHandler{repo: repo}
}

func (h *GetCartHandler) Handle(
	ctx context.Context,
	accountName string,
) (*domain.ShoppingCart, error) {
	return h.repo.Get(ctx, accountName)
}
