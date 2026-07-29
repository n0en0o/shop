package domain

type ShoppingCart struct {
	AccountName string             `json:"account_name" validate:"required,max=20"`
	Items       []ShoppingCartItem `json:"items" validate:"required,min=1,dive"`
}

func (c *ShoppingCart) TotalPrice() float64 {
	var total float64

	for _, item := range c.Items {
		total += float64(item.Quantity) * item.UnitPrice
	}

	return total
}
