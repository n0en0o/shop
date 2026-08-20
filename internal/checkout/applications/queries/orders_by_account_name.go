package queries

import (
	"context"

	"github.com/n0en0o/shop/internal/checkout/applications/interfaces"
	"github.com/n0en0o/shop/internal/checkout/domain"
)

type OrdersByAccountNameQuery struct {
	AccountName string `json:"accountName" validate:"required,max=100"`
}

type OrdersByAccountNameQueryHandler struct {
	repo interfaces.OrderRepository
}

func NewOrdersByAccountNameQueryHandler(
	repo interfaces.OrderRepository,
) *OrdersByAccountNameQueryHandler {
	return &OrdersByAccountNameQueryHandler{
		repo: repo,
	}
}

func (h *OrdersByAccountNameQueryHandler) Handle(
	ctx context.Context,
	q OrdersByAccountNameQuery,
) ([]*domain.Order, error) {
	return h.repo.FindByAccountName(ctx, q.AccountName)
}
