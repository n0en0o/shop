package domain

type ShoppingCart struct {
	AccountName string             `json:"account_name" validate:"required,max=20"`
	Items       []ShoppingCartItem `json:"items" validate:"required,min=1,dive"`
}

func (c *ShoppingCart) TotalPrice() float64 {
	var total float64

	for _, item := range c.Items {
		price := item.UnitPrice
		if item.Discount > 0 || item.FinalPrice > 0 {
			price = item.FinalPrice
		}

		total += float64(item.Quantity) * price
	}

	return total
}
