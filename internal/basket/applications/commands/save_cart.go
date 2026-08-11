package commands

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/n0en0o/marketplace/internal/basket/domain"
	"github.com/n0en0o/marketplace/internal/basket/domain/repositories"
	"github.com/n0en0o/marketplace/internal/promotion/grpc/pb"
)

var validate = validator.New()

type SaveCartCommand struct {
	Cart domain.ShoppingCart `json:"cart" validate:"required"`
}

type SaveCartHandler struct {
	repo        repositories.CartRepository
	promoClient pb.PromotionServiceClient
}

func NewSaveCartHandler(
	repo repositories.CartRepository,
	promoClient pb.PromotionServiceClient,
) *SaveCartHandler {
	return &SaveCartHandler{
		repo:        repo,
		promoClient: promoClient,
	}
}

func (h *SaveCartHandler) Handle(
	ctx context.Context, cmd SaveCartCommand,
) (string, error) {

	if err := validate.Struct(cmd); err != nil {
		return "", err
	}

	if h.promoClient != nil {
		_ = h.applyDiscountsToCart(ctx, &cmd.Cart)
	}

	_, err := h.repo.Save(ctx, cmd.Cart)
	if err != nil {
		return "", err
	}

	return cmd.Cart.AccountName, nil
}

func (h *SaveCartHandler) getDiscountForItem(
	ctx context.Context,
	item *domain.ShoppingCartItem,
) (*pb.GetPromoByCatalogItemResponse, error) {
	req := &pb.GetPromoByCatalogItemRequest{
		CatalogItemId: item.ItemID.String(),
	}

	resp, err := h.promoClient.GetPromoByCatalogItem(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (h *SaveCartHandler) applyDiscountsToCart(
	ctx context.Context,
	cart *domain.ShoppingCart,
) error {
	for i := range cart.Items {
		item := &cart.Items[i]
		disc, err := h.getDiscountForItem(ctx, item)
		if err != nil {
			continue
		}
		val, err := parseDiscountValue(disc)
		if err != nil {
			continue
		}

		item.UnitPrice -= val

	}

	return nil
}

func parseDiscountValue(
	d *pb.GetPromoByCatalogItemResponse,
) (float64, error) {
	value := strings.TrimSpace(d.GetPromotion().GetValue())
	f, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid promo value %q: %v", value, err)
	}

	return f, nil
}
