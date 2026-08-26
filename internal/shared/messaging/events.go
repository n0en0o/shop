package messaging

import (
	"time"

	"github.com/google/uuid"
)

const (
	OrderSubmittedEventType = "OrderSubmittedEvent"
	OrderSubmittedQueue     = "order-submitted-queue"
	OrderSubmittedExchange  = "order-submitted-exchange"
)

type BaseIntegrationEvent struct {
	CorrelationID string    `json:"correlationId"`
	CreationDate  time.Time `json:"creationDate"`
	EventType     string    `json:"eventType"`
}

func NewBaseIntegrationEvent(eventType string) BaseIntegrationEvent {
	return BaseIntegrationEvent{
		CorrelationID: uuid.New().String(),
		CreationDate:  time.Now().UTC(),
		EventType:     eventType,
	}
}

func NewBaseIntegrationEventWithCorrelation(
	eventType, correlationID string,
) BaseIntegrationEvent {
	return BaseIntegrationEvent{
		CorrelationID: correlationID,
		CreationDate:  time.Now().UTC(),
		EventType:     eventType,
	}
}

type OrderItemEventDTO struct {
	CatalogItemName string  `json:"catalogItemName"`
	Quantity        int     `json:"quantity"`
	UnitPrice       float64 `json:"unitPrice"`
}

type OrderSubmittedEvent struct {
	BaseIntegrationEvent

	OrderID     uuid.UUID `json:"orderId"`
	AccountName string    `json:"accountName"`
	TotalPrice  float64   `json:"totalPrice"`

	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`

	Street     string `json:"street"`
	City       string `json:"city"`
	Region     string `json:"region"`
	PostalCode string `json:"postalCode"`

	PaymentMethod int     `json:"paymentMethod"` // 0 = CreditCard, 1 = BankTransfer
	CardName      *string `json:"cardName,omitempty"`
	CardNumber    *string `json:"cardNumber,omitempty"`
	Expiration    *string `json:"expiration,omitempty"`
	CVV           *string `json:"cvv,omitempty"`

	Items []OrderItemEventDTO `json:"items"`
}

func NewOrderSubmittedEvent(correlationID string) *OrderSubmittedEvent {
	return &OrderSubmittedEvent{
		BaseIntegrationEvent: NewBaseIntegrationEventWithCorrelation(
			OrderSubmittedEventType,
			correlationID,
		),
		OrderID: uuid.New(),
		Items:   make([]OrderItemEventDTO, 0),
	}
}
