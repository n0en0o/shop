package domain

type Address struct {
	Street     string `json:"street" validate:"required,max=200"`
	City       string `json:"city" validate:"required,max=100"`
	Region     string `json:"region" validate:"required,max=100"`
	PostalCode string `json:"postalCode" validate:"required,max=20"`
}
