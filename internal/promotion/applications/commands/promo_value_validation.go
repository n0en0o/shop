package commands

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

var (
	ErrInvalidPromoValue  = errors.New("invalid promo value")
	errNegativePromoValue = errors.New("promo value cannot be negative")
)

func validatePromoValue(value string) error {
	amount, err := strconv.ParseFloat(strings.TrimSpace(value), 64)
	if err != nil || math.IsNaN(amount) || math.IsInf(amount, 0) {
		return fmt.Errorf("%w: value must be a valid number", ErrInvalidPromoValue)
	}

	if amount < 0 {
		return fmt.Errorf("%w: %v", ErrInvalidPromoValue, errNegativePromoValue)
	}

	return nil
}
