package domain

type CardDetails struct {
	CardName   string `json:"cardName" validate:"max=100"`
	CardNumber string `json:"cardNumber" validate:"max=20"`
	Expiration string `json:"expiration" validate:"max=10"`
	CVV        string `json:"cvv" validate:"max=10"`
}
