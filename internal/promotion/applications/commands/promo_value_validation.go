package commands

import (
	"errors"
	"fmt"
	"math"
)

var (
	ErrInvalidPromoValue  = errors.New("некорректное значение скидки")
	errNegativePromoValue = errors.New("значение скидки не может быть отрицательным")
)

func validatePromoValue(value float64) error {
	if math.IsNaN(value) || math.IsInf(value, 0) {
		return fmt.Errorf("%w: значение должно быть числом", ErrInvalidPromoValue)
	}

	if value < 0 {
		return fmt.Errorf("%w: %v", ErrInvalidPromoValue, errNegativePromoValue)
	}

	return nil
}
