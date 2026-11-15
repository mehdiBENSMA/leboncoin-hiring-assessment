package money_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.mpi-internal.com/leboncoin/backend-hiring-assessment/common/money"
	"testing"
)

func TestMoneyScaling(t *testing.T) {
	t.Run("convert EUR TO Cents", func(t *testing.T) {
		amountInEur := 156.50
		amountInCents, err := money.ToMinorUnit(amountInEur)

		require.NoError(t, err)
		assert.Equal(t, int64(15650), amountInCents)
	})

	t.Run("convert Cents To EUR", func(t *testing.T) {
		amountInCents := int64(15650)
		amountInEur, err := money.ToMajorUnit(amountInCents)

		require.NoError(t, err)
		assert.Equal(t, float64(156.50), amountInEur)
	})
}
