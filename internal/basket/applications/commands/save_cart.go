package commands

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/n0en0o/marketplace/internal/basket/domain"
	"github.com/n0en0o/marketplace/internal/basket/domain/repositories"
)

var validate = validator.New()

type SaveCartCommand struct {
	Cart domain.ShoppingCart `json:"cart" validate:"required"`
}

type SaveCartHandler struct {
	repo repositories.CartRepository
}

func NewSaveCartHandler(repo repositories.CartRepository) *SaveCartHandler {
	return &SaveCartHandler{repo: repo}
}

func (h *SaveCartHandler) Handle(
	ctx context.Context, cmd SaveCartCommand,
) (string, error) {

	if err := validate.Struct(cmd); err != nil {
		return "", err
	}

	_, err := h.repo.Save(ctx, cmd.Cart)
	if err != nil {
		return "", err
	}

	return cmd.Cart.AccountName, nil
}
