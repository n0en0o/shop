package domain

import "github.com/google/uuid"

type ShoppingCartItem struct {
	ItemID    uuid.UUID `json:"item_id" validate:"required"`
	Quantity  int       `json:"quantity" validate:"required,gt=0"`
	UnitPrice float64   `json:"unit_price" validate:"required,gt=0"`
	ItemTitle *string   `json:"item_title" validate:"required"`
	ItemNote  *string   `json:"item_note"`
}
