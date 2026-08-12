package commands

import (
	"context"

	"github.com/n0en0o/shop/internal/promotion/domain"
	"github.com/n0en0o/shop/internal/promotion/domain/repositories"
)

type UpdatePromoCommand struct {
	ID    string
	Title string
	Value float64
}

type UpdatePromoResult struct {
	Success     bool
	Description string
}

type UpdatePromoHandler struct {
	repo repositories.PromotionRepository
}

func NewUpdatePromoHandler(
	repo repositories.PromotionRepository,
) *UpdatePromoHandler {
	return &UpdatePromoHandler{repo: repo}
}

func (h *UpdatePromoHandler) Handle(
	ctx context.Context,
	cmd UpdatePromoCommand,
) (*UpdatePromoResult, error) {
	if err := validatePromoValue(cmd.Value); err != nil {
		return nil, err
	}

	promo := &domain.Promo{
		ID:    cmd.ID,
		Title: cmd.Title,
		Value: cmd.Value,
	}

	success, err := h.repo.Update(ctx, promo)
	if err != nil {
		return nil, err
	}

	desc := "Акция обновлена"
	if !success {
		desc = "Акция не найдена"
	}

	return &UpdatePromoResult{
		Success:     success,
		Description: desc,
	}, nil
}
