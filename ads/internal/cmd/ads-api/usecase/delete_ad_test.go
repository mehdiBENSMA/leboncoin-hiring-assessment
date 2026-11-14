package usecase_test

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.mpi-internal.com/leboncoin/backend-hiring-assessment/ads/internal/cmd/ads-api/usecase"
	"testing"

	adsdao "github.mpi-internal.com/leboncoin/backend-hiring-assessment/ads/internal/cmd/ads-api/dao"
	"github.mpi-internal.com/leboncoin/backend-hiring-assessment/ads/internal/cmd/ads-api/dao/mock"
)

func TestDeleteAd(t *testing.T) {
	ctx := context.Background()

	ucReq := usecase.DeleteAdRequest{
		AdID:    1,
		OwnerID: "user-1",
	}
	daoReq := adsdao.DeleteAdRequest{
		AdID:    1,
		OwnerID: "user-1",
	}

	t.Run("nominal", func(t *testing.T) {
		m := mock.New(t)

		m.On("DeleteAd", ctx, daoReq).Return(nil)

		err := usecase.NewDeleteAdFunc(m)(ctx, ucReq)

		require.NoError(t, err)
	})

	t.Run("error from dao", func(t *testing.T) {
		m := mock.New(t)
		dbErr := errors.New("connection reset")
		m.On("DeleteAd", ctx, daoReq).Return(dbErr)

		err := usecase.NewDeleteAdFunc(m)(ctx, ucReq)

		require.Error(t, err)
		assert.ErrorIs(t, err, dbErr)
	})
}
