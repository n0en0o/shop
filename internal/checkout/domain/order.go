package domain

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID             uuid.UUID  `json:"id"`
	CreatedBy      *string    `json:"createdBy,omitempty"`
	CreatedAt      time.Time  `json:"createdAt"`
	LastModifiedBy *string    `json:"lastModifiedBy,omitempty"`
	LastModifiedAt *time.Time `json:"lastModifiedAt,omitempty"`

	AccountName          string        `json:"accountName" validate:"required,max=100"`
	TotalAmount          float64       `json:"totalAmount"`
	Items                []OrderItem   `json:"items" validate:"required,min=1,dive"`
	CurrentOrderStatus   OrderStatus   `json:"currentOrderStatus"`
	ContactInfo          Contact       `json:"contactInfo" validate:"required"`
	DeliveryAddress      Address       `json:"deliveryAddress" validate:"required"`
	CurrentPaymentMethod PaymentMethod `json:"currentPaymentMethod"`
	CardDetails          *CardDetails  `json:"cardDetails,omitempty"`
	CurrentPaymentStatus PaymentStatus `json:"currentPaymentStatus"`
}

func NewOrder() *Order {
	return &Order{
		ID:                   uuid.New(),
		CreatedAt:            time.Now().UTC(),
		CurrentOrderStatus:   OrderStatusDraft,
		CurrentPaymentStatus: PaymentStatusPending,
		Items:                make([]OrderItem, 0),
	}
}

func (o *Order) SetCreatedAudit(createdBy string) {
	o.CreatedBy = &createdBy
	o.CreatedAt = time.Now().UTC()
}

func (o *Order) SetModifiedAudit(modifiedBy string) {
	o.LastModifiedBy = &modifiedBy
	now := time.Now().UTC()
	o.LastModifiedAt = &now
}

func (o *Order) CalculateTotalAmount() {
	var total float64
	for _, item := range o.Items {
		total += item.Total()
	}
	o.TotalAmount = total
}
