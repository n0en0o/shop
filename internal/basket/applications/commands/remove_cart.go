package commands

import (
	"context"

	"github.com/n0en0o/marketplace/internal/basket/applications/interfaces"
)

type RemoveCartHandler struct {
	repo interfaces.CartRepository
}

func NewRemoveCartHandler(
	repo interfaces.CartRepository,
) *RemoveCartHandler {
	return &RemoveCartHandler{repo: repo}
}

func (h *RemoveCartHandler) Handle(
	ctx context.Context, accountName string) (bool, error) {
	return h.repo.Remove(ctx, accountName)
}
