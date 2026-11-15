package usecase

import (
	"context"
	"fmt"
	adsdao "github.mpi-internal.com/leboncoin/backend-hiring-assessment/ads/internal/cmd/ads-api/dao"
	"github.mpi-internal.com/leboncoin/backend-hiring-assessment/ads/internal/cmd/ads-api/domain"
	"github.mpi-internal.com/leboncoin/backend-hiring-assessment/common/model"
)

type (
	//ImportAdsFunc import ads owned by the calling user
	ImportAdsFunc func(ctx context.Context, req ImportAdsRequest) ([]domain.AdID, error)

	ImportAdsItemRequest struct {
		Title    string
		Price    float64
		PhotoURL string
	}

	ImportAdsRequest struct {
		Ads     []ImportAdsItemRequest
		OwnerId model.UserID
	}

	importAds struct {
		adDAO adsdao.DAO
	}
)

func NewImportAdsFunc(dDAO adsdao.DAO) ImportAdsFunc {
	return importAds{adDAO: dDAO}.importAds
}

func (uc importAds) importAds(ctx context.Context, req ImportAdsRequest) ([]domain.AdID, error) {
	tx, err := uc.adDAO.Begin(ctx)

	ids := make([]domain.AdID, 0, len(req.Ads))

	if err != nil {
		return ids, err
	}

	for _, ad := range req.Ads {
		adID, err := uc.adDAO.CreateAd(ctx, adsdao.CreateAdRequest{
			Title:    ad.Title,
			Price:    ad.Price,
			PhotoURL: ad.PhotoURL,
			OwnerID:  req.OwnerId,
		})

		if err != nil {
			err = tx.Rollback()
			if err != nil {
				return nil, fmt.Errorf("failed to rollback transaction: %w", err)
			}
			return nil, fmt.Errorf("failed to import ad (%q) : %w", ad.Title, err)
		}

		ids = append(ids, adID)
	}

	err = tx.Commit()
	if err != nil {
		return ids, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return ids, nil
}
