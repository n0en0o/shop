package commands

import (
	"context"

	"github.com/n0en0o/marketplace/internal/basket/applications/interfaces"
	"github.com/n0en0o/marketplace/internal/basket/domain"
)

type SaveCartCommand struct {
	Cart domain.ShoppingCart `json:"cart"`
}

type SaveCartHandler struct {
	repo interfaces.CartRepository
}

func NewSaveCartHandler(repo interfaces.CartRepository) *SaveCartHandler {
	return &SaveCartHandler{repo: repo}
}

func (h *SaveCartHandler) Handle(
	ctx context.Context, cmd SaveCartCommand,
) (string, error) {
	_, err := h.repo.Save(ctx, cmd.Cart)
	if err != nil {
		return "", err
	}

	return cmd.Cart.AccountName, nil
}
