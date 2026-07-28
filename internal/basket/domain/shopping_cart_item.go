package domain

import "github.com/google/uuid"

type ShoppingCartItem struct {
	ItemID    uuid.UUID `json:"item_id"`
	Quantity  int       `json:"quantity"`
	UnitPrice float64   `json:"unit_price"`
	ItemTitle *string   `json:"item_title"`
	ItemNote  *string   `json:"item_note"`
}
