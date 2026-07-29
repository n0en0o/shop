package shared

import "fmt"

type NotFoundError struct {
	Resource string
	Key      any
}

func NewNotFoundError(resource string, key any) *NotFoundError {
	return &NotFoundError{Resource: resource, Key: key}
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("объект %q с ключом %v не найден", e.Resource, e.Key)
}
