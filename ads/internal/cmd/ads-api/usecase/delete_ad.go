package usecase

import (
	"context"
	"errors"
	"fmt"
	adsdao "github.mpi-internal.com/leboncoin/backend-hiring-assessment/ads/internal/cmd/ads-api/dao"
	"github.mpi-internal.com/leboncoin/backend-hiring-assessment/ads/internal/cmd/ads-api/domain"
)

type (
	//DeleteAdFunc delete an ad owned by the calling user
	DeleteAdFunc func(ctx context.Context, req DeleteAdRequest) error

	DeleteAdRequest struct {
		AdID    domain.AdID
		OwnerID domain.UserID
	}

	deleteAd struct {
		adDAO adsdao.DAO
	}
)

func NewDeleteAdFunc(dDAO adsdao.DAO) DeleteAdFunc {
	return deleteAd{adDAO: dDAO}.deleteAd
}

func (uc deleteAd) deleteAd(ctx context.Context, req DeleteAdRequest) error {
	err := uc.adDAO.DeleteAd(ctx, adsdao.DeleteAdRequest{
		AdID:    req.AdID,
		OwnerID: req.OwnerID,
	})

	if errors.Is(err, adsdao.ErrNotFound) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("unable to delete ad: %w", err)
	}

	return nil
}
