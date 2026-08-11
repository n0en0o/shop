package commands

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/n0en0o/marketplace/internal/promotion/applications/interfaces"
	"github.com/n0en0o/marketplace/internal/promotion/domain"
)

type CreatePromoCommand struct {
	CatalogItemID string
	Title         string
	Value         string
}

type CreatePromoResult struct {
	ID          string
	Success     bool
	Description string
}

type CreatePromoHandler struct {
	repo interfaces.PromotionRepository
}

func NewCreatePromoHandler(
	repo interfaces.PromotionRepository,
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

	description := "Failed to create promotion"
	if success {
		description = "Promotion created successfully"
	}

	return &CreatePromoResult{
		ID:          promo.ID,
		Success:     success,
		Description: description,
	}, nil
}
