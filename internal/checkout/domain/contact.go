package domain

type Contact struct {
	FirstName string `json:"firstName" validate:"required,max=100"`
	LastName  string `json:"lastName" validate:"required,max=100"`
	Email     string `json:"email" validate:"required,email,max=255"`
}
