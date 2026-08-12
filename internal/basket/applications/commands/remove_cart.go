package commands

import (
	"context"

	"github.com/n0en0o/shop/internal/basket/domain/repositories"
)

type RemoveCartHandler struct {
	repo repositories.CartRepository
}

func NewRemoveCartHandler(
	repo repositories.CartRepository,
) *RemoveCartHandler {
	return &RemoveCartHandler{repo: repo}
}

func (h *RemoveCartHandler) Handle(
	ctx context.Context, accountName string) (bool, error) {
	return h.repo.Remove(ctx, accountName)
}
