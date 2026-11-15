package usecase

import (
	"context"
	"errors"
	"fmt"
	adsdao "github.mpi-internal.com/leboncoin/backend-hiring-assessment/ads/internal/cmd/ads-api/dao"
	"github.mpi-internal.com/leboncoin/backend-hiring-assessment/ads/internal/cmd/ads-api/domain"
	"time"
)

const MAX_ADS_PER_DAY = 5

var (
	TooManyAdErr = errors.New("too many created today, comeback tomorrow")
)

type (
	// CreateAdFunc creates an ad owned by the calling user.
	CreateAdFunc func(ctx context.Context, req CreateAdRequest) (domain.AdID, error)

	// CreateAdRequest describes the ad to create. OwnerID is the authenticated user: a caller can never create an ad on behalf of someone else.
	CreateAdRequest struct {
		Title    string
		Price    float64
		PhotoURL string
		OwnerID  domain.UserID
	}

	createAd struct {
		adDAO adsdao.DAO
	}
)

// NewCreateAdFunc returns a CreateAdFunc backed by the given DAO.
func NewCreateAdFunc(adDAO adsdao.DAO) CreateAdFunc {
	return createAd{adDAO: adDAO}.createAd
}

func (uc createAd) createAd(ctx context.Context, req CreateAdRequest) (domain.AdID, error) {

	now := time.Now()
	todayAtMidnight := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		0, 0, 0, 0,
		now.Location(),
	)
	search := adsdao.SearchAdsRequest{
		CreatedAfter:  todayAtMidnight,
		CreatedBefore: now,
	}

	existingAds, err := uc.adDAO.SearchAds(ctx, search)

	if err != nil {
		return 0, fmt.Errorf("searching ads: %w", err)
	}

	if len(existingAds) > MAX_ADS_PER_DAY {
		return 0, TooManyAdErr
	}

	adID, err := uc.adDAO.CreateAd(ctx, adsdao.CreateAdRequest{
		Title:    req.Title,
		Price:    req.Price,
		PhotoURL: req.PhotoURL,
		OwnerID:  req.OwnerID,
	})
	if err != nil {
		return 0, fmt.Errorf("unable to create ad: %w", err)
	}

	return adID, nil
}
