package commands

import (
	"context"
	"errors"
	"math"

	"github.com/go-playground/validator/v10"
	"github.com/n0en0o/shop/internal/basket/domain"
	"github.com/n0en0o/shop/internal/basket/domain/repositories"
	"github.com/n0en0o/shop/internal/promotion/grpc/pb"
)

var validate = validator.New()

var errInvalidPromoValue = errors.New("некорректное значение скидки")

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

	_ = h.applyDiscountsToCart(ctx, &cmd.Cart)

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
		item.Discount = 0
		item.FinalPrice = item.UnitPrice

		if h.promoClient == nil {
			continue
		}

		disc, err := h.getDiscountForItem(ctx, item)
		if err != nil {
			continue
		}
		val, err := discountValue(disc)
		if err != nil {
			continue
		}

		if val < 0 {
			continue
		}
		if val > item.UnitPrice {
			val = item.UnitPrice
		}

		item.Discount = val
		item.FinalPrice = item.UnitPrice - val
	}

	return nil
}

func discountValue(
	d *pb.GetPromoByCatalogItemResponse,
) (float64, error) {
	value := d.GetPromotion().GetValue()
	if math.IsNaN(value) || math.IsInf(value, 0) || value < 0 {
		return 0, errInvalidPromoValue
	}

	return value, nil
}
