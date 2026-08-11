package commands

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/n0en0o/marketplace/internal/promotion/domain"
	"github.com/n0en0o/marketplace/internal/promotion/domain/repositories"
)

type CreatePromoCommand struct {
	CatalogItemID string
	Title         string
	Value         float64
}

type CreatePromoResult struct {
	ID          string
	Success     bool
	Description string
}

type CreatePromoHandler struct {
	repo repositories.PromotionRepository
}

func NewCreatePromoHandler(
	repo repositories.PromotionRepository,
) *CreatePromoHandler {
	return &CreatePromoHandler{
		repo: repo,
	}
}

func (h *CreatePromoHandler) Handle(
	ctx context.Context,
	cmd CreatePromoCommand,
) (*CreatePromoResult, error) {
	if err := validatePromoValue(cmd.Value); err != nil {
		return nil, err
	}

	promo := &domain.Promo{
		ID:            uuid.New().String(),
		CatalogItemID: cmd.CatalogItemID,
		Title:         cmd.Title,
		Value:         cmd.Value,
	}

	success, err := h.repo.Create(ctx, promo)
	if err != nil {
		if errors.Is(err, domain.ErrPromoAlreadyExists) {
			return nil, err
		}
		return nil, fmt.Errorf("create promotion: %w", err)
	}

	description := "Акция не создана"
	if success {
		description = "Акция создана успешно"
	}

	return &CreatePromoResult{
		ID:          promo.ID,
		Success:     success,
		Description: description,
	}, nil
}
