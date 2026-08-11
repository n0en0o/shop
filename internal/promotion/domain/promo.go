package domain

import "errors"

var ErrPromoAlreadyExists = errors.New("promotion already exists")

type Promo struct {
	ID            string
	CatalogItemID string
	Title         string
	Value         string
}
