package domain

type OrderStatus string

type OrdersStatus = OrderStatus

const (
	OrderStatusDraft     OrderStatus = "Draft"
	OrderStatusSubmitted OrderStatus = "Submitted"
	OrderStatusPaid      OrderStatus = "Paid"
	OrderStatusShipped   OrderStatus = "Shipped"
	OrderStatusCancelled OrderStatus = "Cancelled"
)

func (s OrderStatus) String() string {
	return string(s)
}

func (s OrderStatus) IsValid() bool {
	switch s {
	case OrderStatusDraft,
		OrderStatusSubmitted,
		OrderStatusPaid,
		OrderStatusShipped,
		OrderStatusCancelled:
		return true
	}

	return false
}
