package commands

import (
	"context"
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/n0en0o/shop/internal/basket/domain/repositories"
	"github.com/n0en0o/shop/internal/shared/messaging"
)

type CheckoutCartRequest struct {
	AccountName   string  `json:"accountName" validate:"required,min=3,max=20"`
	FirstName     string  `json:"firstName" validate:"required,max=50"`
	LastName      string  `json:"lastName" validate:"required,max=50"`
	Email         string  `json:"email" validate:"required,email"`
	Street        string  `json:"street" validate:"required,max=200"`
	City          string  `json:"city" validate:"required,max=100"`
	Region        string  `json:"region" validate:"required"`
	PostalCode    string  `json:"postalCode" validate:"required"`
	PaymentMethod int     `json:"paymentMethod"` // 0 = CreditCard, 1 = BankTransfer
	CardName      *string `json:"cardName,omitempty"`
	CardNumber    *string `json:"cardNumber,omitempty"`
	Expiration    *string `json:"expiration,omitempty"`
	CVV           *string `json:"cvv,omitempty"`
}

type CheckoutCartCommand struct {
	CheckoutCartRequest
	CorrelationID string
}

type CheckoutCartResult struct {
	OrderID       uuid.UUID `json:"orderId"`
	CorrelationID string    `json:"correlationId"`
	CartRemoved   bool      `json:"cartRemoved"`
}

type CheckoutCartResponse struct {
	OrderID       uuid.UUID `json:"orderId"`
	CorrelationID string    `json:"correlationId"`
	Message       string    `json:"message"`
}

type CheckoutCartHandler struct {
	cartRepo  repositories.CartRepository
	publisher *messaging.RabbitMQPublisher
	validate  *validator.Validate
}

func NewCheckoutCartHandler(
	cartRepo repositories.CartRepository,
	publisher *messaging.RabbitMQPublisher,
) *CheckoutCartHandler {
	return &CheckoutCartHandler{
		cartRepo:  cartRepo,
		publisher: publisher,
		validate:  validator.New(),
	}
}

func (h *CheckoutCartHandler) Handle(
	ctx context.Context,
	cmd CheckoutCartCommand,
) (*CheckoutCartResult, error) {
	if err := h.validate.Struct(cmd.CheckoutCartRequest); err != nil {
		return nil, err
	}

	if cmd.PaymentMethod == 0 {
		if cmd.CardNumber == nil || *cmd.CardNumber == "" {
			return nil, errors.New("cardNumber обязателен для CreditCard")
		}
		if cmd.CVV == nil || *cmd.CVV == "" {
			return nil, errors.New("CVV обязателен для CreditCard")
		}
	}

	cart, err := h.cartRepo.Get(ctx, cmd.AccountName)
	if err != nil {
		return nil, err
	}
	if cart == nil {
		return nil, errors.New("корзина не найдена")
	}
	if len(cart.Items) == 0 {
		return nil, errors.New("корзина пуста")
	}

	event := messaging.NewOrderSubmittedEvent(cmd.CorrelationID)
	event.AccountName = cmd.AccountName
	event.TotalPrice = cart.TotalPrice()
	event.FirstName = cmd.FirstName
	event.LastName = cmd.LastName
	event.Email = cmd.Email
	event.Street = cmd.Street
	event.City = cmd.City
	event.Region = cmd.Region
	event.PostalCode = cmd.PostalCode
	event.PaymentMethod = cmd.PaymentMethod
	event.CardName = cmd.CardName
	event.CardNumber = cmd.CardNumber
	event.Expiration = cmd.Expiration
	event.CVV = cmd.CVV

	for _, item := range cart.Items {
		catalogItemName := ""
		if item.ItemTitle != nil {
			catalogItemName = *item.ItemTitle
		}

		event.Items = append(event.Items, messaging.OrderItemEventDTO{
			CatalogItemName: catalogItemName,
			Quantity:        item.Quantity,
			UnitPrice:       item.UnitPrice,
		})
	}

	if h.publisher == nil {
		return nil, errors.New("RabbitMQ publisher is not configured")
	}

	if err := h.publisher.PublishOrderSubmitted(ctx, event); err != nil {
		return nil, err
	}

	return &CheckoutCartResult{
		OrderID:       event.OrderID,
		CorrelationID: cmd.CorrelationID,
		CartRemoved:   false,
	}, nil
}
