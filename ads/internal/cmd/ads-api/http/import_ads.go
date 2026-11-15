package http

import (
	"encoding/json"
	"net/http"

	"github.mpi-internal.com/leboncoin/backend-hiring-assessment/ads/internal/cmd/ads-api/domain"
	"github.mpi-internal.com/leboncoin/backend-hiring-assessment/ads/internal/cmd/ads-api/usecase"
	"github.mpi-internal.com/leboncoin/backend-hiring-assessment/common/middleware"
)

type importAdsBody struct {
	Ads []createAdBody `json:"ads"`
}

type importAdsResponse struct {
	IDs []int64 `json:"ids"`
}

func handlerImportAds(importAds usecase.ImportAdsFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var body importAdsBody
		if err := decodeBody(r, &body); err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		ads := make([]usecase.ImportAdsItemRequest, 0, len(body.Ads))
		for _, ad := range body.Ads {
			ads = append(ads, usecase.ImportAdsItemRequest{
				Title:    ad.Title,
				Price:    ad.Price,
				PhotoURL: ad.PhotoURL,
			})
		}

		if len(ads) == 0 {
			rw.Header().Set("Content-Type", "application/json")
			rw.WriteHeader(http.StatusCreated)

			_ = json.NewEncoder(rw).Encode(importAdsResponse{
				IDs: []int64{},
			})
			return
		}

		adIDS, err := importAds(r.Context(), usecase.ImportAdsRequest{
			Ads:     ads,
			OwnerId: domain.UserID(middleware.UserIDFromRequest(r)),
		})

		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusCreated)

		ids := make([]int64, 0, len(adIDS))

		for _, id := range adIDS {
			ids = append(ids, int64(id))
		}
		_ = json.NewEncoder(rw).Encode(importAdsResponse{IDs: ids})
	}
}
