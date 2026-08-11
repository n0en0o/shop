package commands

import (
	"context"

	"github.com/n0en0o/marketplace/internal/promotion/applications/interfaces"
	"github.com/n0en0o/marketplace/internal/promotion/domain"
)

type UpdatePromoCommand struct {
	ID    string
	Title string
	Value string
}

type UpdatePromoResult struct {
	Success     bool
	Description string
}

type UpdatePromoHandler struct {
	repo interfaces.PromotionRepository
}

func NewUpdatePromoHandler(
	repo interfaces.PromotionRepository,
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

	desc := "обновлено"
	if !success {
		desc = "ошибка при обновления"
	}

	return &UpdatePromoResult{
		Success:     success,
		Description: desc,
	}, nil
}
