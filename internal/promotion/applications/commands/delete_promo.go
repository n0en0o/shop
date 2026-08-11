package commands

import (
	"context"

	"github.com/n0en0o/marketplace/internal/promotion/domain/repositories"
)

type DeletePromoCommand struct {
	ID string
}

type DeletePromoResult struct {
	Success     bool
	Description string
}

type DeletePromoHandler struct {
	repo repositories.PromotionRepository
}

func NewDeletePromoHandler(
	repo repositories.PromotionRepository,
) *DeletePromoHandler {
	return &DeletePromoHandler{repo: repo}
}

func (h *DeletePromoHandler) Handle(
	ctx context.Context,
	cmd DeletePromoCommand,
) (*DeletePromoResult, error) {
	success, err := h.repo.Delete(ctx, cmd.ID)
	if err != nil {
		return nil, err
	}

	desc := "Акция удалена"
	if !success {
		desc = "Акция не найдена"
	}

	return &DeletePromoResult{
		Success:     success,
		Description: desc,
	}, nil
}
