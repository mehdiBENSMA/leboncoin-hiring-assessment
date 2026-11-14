package http

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.mpi-internal.com/leboncoin/backend-hiring-assessment/ads/internal/cmd/ads-api/usecase"
	"github.mpi-internal.com/leboncoin/backend-hiring-assessment/common/middleware"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newDeleteAdServer(fn usecase.DeleteAdFunc) *httptest.Server {
	mux := http.NewServeMux()
	mux.Handle("DELETE /ads/{adID}", middleware.RequireAuth(handleDeleteAd(fn)))

	return httptest.NewServer(mux)
}

func TestHandleDeleteAd(t *testing.T) {
	t.Run("return 204 when deleting the ad", func(t *testing.T) {

		srv := newDeleteAdServer(func(_ context.Context, req usecase.DeleteAdRequest) error {
			return nil
		})

		defer srv.Close()
		ctx := context.Background()
		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodDelete,
			srv.URL+"/ads/42",
			nil,
		)
		req.Header.Set("Authorization", "user-1")

		if err != nil {
			panic(err)
		}

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("return 401 when Authorization header is missing", func(t *testing.T) {
		srv := newDeleteAdServer(func(_ context.Context, req usecase.DeleteAdRequest) error {
			t.Fatal("usecase should not be called")
			return nil
		})

		defer srv.Close()
		ctx := context.Background()
		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodDelete,
			srv.URL+"/ads/42",
			nil,
		)

		if err != nil {
			panic(err)
		}

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}
