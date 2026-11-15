package money

import (
	"fmt"
	"github.com/govalues/money"
)

const defaultCurrency = "EUR"
const defaultScale = 2
const defaultWhole = 0

func ToMinorUnit(price float64) (int64, error) {
	amt, err := money.NewAmountFromFloat64(defaultCurrency, price)
	if err != nil {
		return 0, fmt.Errorf("unable to create minor unit amount: %w", err)
	}

	whole, frac, ok := amt.Int64(defaultScale)
	if !ok {
		return 0, fmt.Errorf("unable to create minor unit amount")
	}

	cents := whole*100 + frac
	return cents, nil
}

func ToMajorUnit(priceInCent int64) (float64, error) {
	amt, err := money.NewAmount(defaultCurrency, priceInCent, defaultScale)
	if err != nil {
		return 0, fmt.Errorf("unable to create minor unit amount: %w", err)
	}

	amountInMajor, ok := amt.Float64()
	if !ok {
		return 0, fmt.Errorf("unable to create minor unit amount")
	}
	return amountInMajor, nil
}
