package entities

import "github.com/google/uuid"

type BaseEntity struct {
	ID    uuid.UUID `json:"id"`
	Title string    `json:"title"`
}
