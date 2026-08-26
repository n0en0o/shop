package messaging

import (
	"context"
	"log"

	"github.com/n0en0o/shop/internal/checkout/applications/commands"
	sharedmessaging "github.com/n0en0o/shop/internal/shared/messaging"
)

type OrderSubmittedConsumer struct {
	handler *commands.ProcessOrderSubmissionHandler
}

func NewOrderSubmittedConsumer(
	handler *commands.ProcessOrderSubmissionHandler,
) *OrderSubmittedConsumer {
	return &OrderSubmittedConsumer{handler: handler}
}

func (c *OrderSubmittedConsumer) HandleMessage(
	ctx context.Context,
	body []byte,
) error {
	event, err := sharedmessaging.ParseMessage[sharedmessaging.OrderSubmittedEvent](body)
	if err != nil {
		log.Printf("[order-consumer] ошибка парсинга: %v", err)
		return err
	}

	log.Printf(
		"[order-consumer] получен OrderSubmittedEvent: orderId=%s, account=%s",
		event.OrderID,
		event.AccountName,
	)

	cmd := c.mapEventToCommand(event)

	result, err := c.handler.Handle(ctx, cmd)
	if err != nil {
		log.Printf("[order-consumer] ошибка обработки: %v", err)
		return err
	}

	log.Printf("[order-consumer] заказ обработан: %s", result.OrderID)
	return nil
}

func (c *OrderSubmittedConsumer) mapEventToCommand(
	event *sharedmessaging.OrderSubmittedEvent,
) commands.ProcessOrderSubmissionCommand {
	items := make([]commands.OrderItemDTO, len(event.Items))
	for i, item := range event.Items {
		items[i] = commands.OrderItemDTO{
			CatalogItemName: item.CatalogItemName,
			Quantity:        item.Quantity,
			UnitPrice:       item.UnitPrice,
		}
	}

	return commands.ProcessOrderSubmissionCommand{
		OrderID:       event.OrderID,
		AccountName:   event.AccountName,
		TotalPrice:    event.TotalPrice,
		FirstName:     event.FirstName,
		LastName:      event.LastName,
		Email:         event.Email,
		Street:        event.Street,
		City:          event.City,
		Region:        event.Region,
		PostalCode:    event.PostalCode,
		PaymentMethod: event.PaymentMethod,
		CardName:      event.CardName,
		CardNumber:    event.CardNumber,
		Expiration:    event.Expiration,
		CVV:           event.CVV,
		CorrelationID: event.CorrelationID,
		Items:         items,
	}
}
