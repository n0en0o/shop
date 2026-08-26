package commands

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"github.com/n0en0o/shop/internal/checkout/applications/interfaces"
	"github.com/n0en0o/shop/internal/checkout/domain"
	"github.com/n0en0o/shop/internal/shared"
)

type ProcessOrderSubmissionCommand struct {
	OrderID     uuid.UUID `validate:"required"`
	AccountName string    `validate:"required,min=3,max=100"`
	TotalPrice  float64   `validate:"gt=0"`

	FirstName string `validate:"required"`
	LastName  string `validate:"required"`
	Email     string `validate:"required,email"`

	Street     string `validate:"required"`
	City       string `validate:"required"`
	Region     string `validate:"required"`
	PostalCode string `validate:"required"`

	PaymentMethod int

	CardName   *string
	CardNumber *string
	Expiration *string
	CVV        *string

	CorrelationID string

	Items []OrderItemDTO `validate:"required,min=1,dive"`
}

type OrderItemDTO struct {
	CatalogItemName string  `validate:"required"`
	Quantity        int     `validate:"gt=0"`
	UnitPrice       float64 `validate:"gt=0"`
}

type ProcessOrderSubmissionResult struct {
	OrderID uuid.UUID `json:"orderId"`
}

type ProcessOrderSubmissionHandler struct {
	repo      interfaces.OrderRepository
	validator *validator.Validate
}

func NewProcessOrderSubmissionHandler(
	repo interfaces.OrderRepository,
) *ProcessOrderSubmissionHandler {
	return &ProcessOrderSubmissionHandler{
		repo:      repo,
		validator: validator.New(),
	}
}

func (h *ProcessOrderSubmissionHandler) Handle(
	ctx context.Context,
	cmd ProcessOrderSubmissionCommand,
) (*ProcessOrderSubmissionResult, error) {
	if err := h.validator.Struct(cmd); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	existingOrder, err := h.repo.Get(ctx, cmd.OrderID)
	if err != nil {
		var notFound *shared.NotFoundError
		if !errors.As(err, &notFound) {
			return nil, fmt.Errorf("не удалось проверить существующий заказ: %w", err)
		}
	}

	if existingOrder != nil {
		log.Printf("[process-order] заказ %s уже существует", cmd.OrderID)
		return &ProcessOrderSubmissionResult{OrderID: existingOrder.ID}, nil
	}

	order := h.mapToOrder(cmd)

	if err := h.repo.Create(ctx, order); err != nil {
		return nil, fmt.Errorf("не удалось создать заказ: %w", err)
	}

	log.Printf("[process-order] создан заказ %s для %s", order.ID, order.AccountName)

	return &ProcessOrderSubmissionResult{OrderID: order.ID}, nil
}

func paymentMethodFromInt(method int) domain.PaymentMethod {
	switch method {
	case 0:
		return domain.PaymentMethodCreditCard
	case 1:
		return domain.PaymentMethodBankTransfer
	default:
		return ""
	}
}

func (h *ProcessOrderSubmissionHandler) mapToOrder(
	cmd ProcessOrderSubmissionCommand,
) *domain.Order {
	order := domain.NewOrder()
	order.ID = cmd.OrderID
	order.AccountName = cmd.AccountName
	order.TotalAmount = cmd.TotalPrice

	order.ContactInfo = domain.Contact{
		FirstName: cmd.FirstName,
		LastName:  cmd.LastName,
		Email:     cmd.Email,
	}

	order.DeliveryAddress = domain.Address{
		Street:     cmd.Street,
		City:       cmd.City,
		Region:     cmd.Region,
		PostalCode: cmd.PostalCode,
	}

	order.CurrentPaymentMethod = paymentMethodFromInt(cmd.PaymentMethod)
	if cmd.CardName != nil && cmd.CardNumber != nil {
		order.CardDetails = &domain.CardDetails{
			CardName:   *cmd.CardName,
			CardNumber: *cmd.CardNumber,
			Expiration: stringOrEmpty(cmd.Expiration),
			CVV:        stringOrEmpty(cmd.CVV),
		}
	}

	order.Items = make([]domain.OrderItem, len(cmd.Items))
	for i, item := range cmd.Items {
		order.Items[i] = domain.OrderItem{
			CatalogItemName: item.CatalogItemName,
			Quantity:        item.Quantity,
			UnitPrice:       item.UnitPrice,
		}
	}

	order.CurrentOrderStatus = domain.OrderStatusSubmitted

	return order
}

func stringOrEmpty(s *string) string {
	if s == nil {
		return ""
	}

	return *s
}
