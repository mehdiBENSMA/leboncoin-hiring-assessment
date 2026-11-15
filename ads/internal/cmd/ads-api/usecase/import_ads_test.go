package usecase_test

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	adsdao "github.mpi-internal.com/leboncoin/backend-hiring-assessment/ads/internal/cmd/ads-api/dao"
	"github.mpi-internal.com/leboncoin/backend-hiring-assessment/ads/internal/cmd/ads-api/dao/mock"
	"github.mpi-internal.com/leboncoin/backend-hiring-assessment/ads/internal/cmd/ads-api/domain"
	"github.mpi-internal.com/leboncoin/backend-hiring-assessment/ads/internal/cmd/ads-api/usecase"
	"testing"
)

func TestImportAds(t *testing.T) {
	ctx := context.Background()

	importAd1Request := usecase.ImportAdsItemRequest{
		Title:    "Mountain bike",
		Price:    450.00,
		PhotoURL: "https://example.com/mountain-bike.jpg",
	}

	importAd2Request := usecase.ImportAdsItemRequest{
		Title:    "Wireless headphones",
		Price:    89.99,
		PhotoURL: "https://example.com/headphones.jpg",
	}

	importAd3Request := usecase.ImportAdsItemRequest{
		Title:    "Gaming keyboard",
		Price:    129.50,
		PhotoURL: "https://example.com/keyboard.jpg",
	}

	importAd4Request := usecase.ImportAdsItemRequest{
		Title:    "Office chair",
		Price:    275.00,
		PhotoURL: "https://example.com/office-chair.jpg",
	}

	importAdsRequest := make([]usecase.ImportAdsItemRequest, 0, 4)
	importAdsRequest = append(importAdsRequest, importAd1Request, importAd2Request, importAd3Request, importAd4Request)

	ucReq := usecase.ImportAdsRequest{
		Ads:     importAdsRequest,
		OwnerId: "user-1",
	}

	daoReq1 := adsdao.CreateAdRequest{
		Title:    "Mountain bike",
		Price:    450.00,
		PhotoURL: "https://example.com/mountain-bike.jpg",
		OwnerID:  "user-1",
	}

	daoReq2 := adsdao.CreateAdRequest{
		Title:    "Wireless headphones",
		Price:    89.99,
		PhotoURL: "https://example.com/headphones.jpg",
		OwnerID:  "user-1",
	}

	daoReq3 := adsdao.CreateAdRequest{
		Title:    "Gaming keyboard",
		Price:    129.50,
		PhotoURL: "https://example.com/keyboard.jpg",
		OwnerID:  "user-1",
	}

	daoReq4 := adsdao.CreateAdRequest{
		Title:    "Office chair",
		Price:    275.00,
		PhotoURL: "https://example.com/office-chair.jpg",
		OwnerID:  "user-1",
	}

	t.Run("nominal", func(t *testing.T) {
		m := mock.New(t)
		m.On("Begin", ctx).Return(nil, nil).Once()
		m.On("CreateAd", ctx, daoReq1).Return(domain.AdID(1), nil).Once()
		m.On("CreateAd", ctx, daoReq2).Return(domain.AdID(2), nil).Once()
		m.On("CreateAd", ctx, daoReq3).Return(domain.AdID(3), nil).Once()
		m.On("CreateAd", ctx, daoReq4).Return(domain.AdID(4), nil).Once()

		m.On("Commit").Return(nil).Once()

		got, err := usecase.NewImportAdsFunc(m)(ctx, ucReq)

		require.NoError(t, err)
		assert.Equal(t, len(importAdsRequest), len(got))
	})

	t.Run("call rollback on Error", func(t *testing.T) {
		dbErr := errors.New("connection reset")

		m := mock.New(t)

		m.On("Begin", ctx).Return(nil, nil).Once()
		m.On("CreateAd", ctx, daoReq1).Return(domain.AdID(1), nil).Once()
		m.On("CreateAd", ctx, daoReq2).Return(domain.AdID(2), nil).Once()
		m.On("CreateAd", ctx, daoReq3).Return(domain.AdID(3), nil).Once()
		m.On("CreateAd", ctx, daoReq4).Return(0, dbErr).Once()

		m.On("Rollback").Return(nil).Once()

		_, err := usecase.NewImportAdsFunc(m)(ctx, ucReq)

		require.Error(t, err)
	})
}
