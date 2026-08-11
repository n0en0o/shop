package domain

import "errors"

var ErrPromoAlreadyExists = errors.New("акция уже существует")

type Promo struct {
	ID            string
	CatalogItemID string
	Title         string
	Value         float64
}
