package domain

type OrderItem struct {
	CatalogItemName string  `json:"catalogItemName" validate:"required,max=200"`
	Quantity        int     `json:"quantity" validate:"required,gt=0"`
	UnitPrice       float64 `json:"unitPrice" validate:"required,gte=0"`
}

func (i OrderItem) Total() float64 {
	return float64(i.Quantity) * i.UnitPrice
}
