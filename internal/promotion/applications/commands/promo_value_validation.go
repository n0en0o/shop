package commands

import (
	"errors"
	"math"
	"strconv"
	"strings"
)

var (
	errInvalidPromoValue  = errors.New("value must be a valid number")
	errNegativePromoValue = errors.New("value cannot be negative")
)

func validatePromoValue(value string) error {
	amount, err := strconv.ParseFloat(strings.TrimSpace(value), 64)
	if err != nil || math.IsNaN(amount) || math.IsInf(amount, 0) {
		return errInvalidPromoValue
	}

	if amount < 0 {
		return errNegativePromoValue
	}

	return nil
}
