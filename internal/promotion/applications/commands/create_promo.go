package commands

import (
	"context"
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
		return &CreatePromoResult{
			ID:          "",
			Success:     false,
			Description: err.Error(),
		}, nil
	}

	existing, err := h.repo.FindByCatalogItem(ctx, cmd.CatalogItemID)

	if err != nil {
		return &CreatePromoResult{
			ID:          "",
			Success:     false,
			Description: fmt.Sprintf("FindByCatalogItem error: %v", err),
		}, nil
	}

	if existing != nil {
		msg := fmt.Sprintf(
			"promotion is already exists for %s (ID: %s)",
			cmd.CatalogItemID,
			existing.ID,
		)

		return &CreatePromoResult{
			ID:          existing.ID,
			Success:     false,
			Description: msg,
		}, nil
	}

	promo := &domain.Promo{
		ID:            uuid.New().String(),
		CatalogItemID: cmd.CatalogItemID,
		Title:         cmd.Title,
		Value:         cmd.Value,
	}

	success, err := h.repo.Create(ctx, promo)
	if err != nil {
		return &CreatePromoResult{
			ID:          "",
			Success:     false,
			Description: fmt.Sprintf("Failed to create promotion: %v", err),
		}, nil
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
