package queries

import (
	"context"

	"github.com/google/uuid"

	"github.com/n0en0o/shop/internal/checkout/applications/interfaces"
	"github.com/n0en0o/shop/internal/checkout/domain"
)

type OrderByIDQuery struct {
	ID uuid.UUID `json:"id" validate:"required,uuid"`
}

type OrderByIDQueryHandler struct {
	repo interfaces.OrderRepository
}

func NewOrderByIDQueryHandler(
	repo interfaces.OrderRepository,
) *OrderByIDQueryHandler {
	return &OrderByIDQueryHandler{
		repo: repo,
	}
}

func (h *OrderByIDQueryHandler) Handle(
	ctx context.Context,
	q OrderByIDQuery,
) (*domain.Order, error) {
	return h.repo.Get(ctx, q.ID)
}
