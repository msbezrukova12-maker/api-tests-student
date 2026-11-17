package advertisementCreate

import (
	"image/color"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/tests/fixtures"
)

func TestCreateAdvertisementValidBoundaries(t *testing.T) {
	user := fixtures.Login(t)

	cases := []struct {
		name    string
		prepare func(*models.CreateAdvertisementRequest)
	}{
		{
			name: "title_50_characters",
			prepare: func(request *models.CreateAdvertisementRequest) {
				request.Fields["title"] = strings.Repeat("a", 50)
			},
		},
		{
			name: "description_500_characters",
			prepare: func(request *models.CreateAdvertisementRequest) {
				request.Fields["description"] = strings.Repeat("b", 500)
			},
		},
		{
			name: "price_0",
			prepare: func(request *models.CreateAdvertisementRequest) {
				request.Fields["price"] = "0"
			},
		},
		{
			name: "price_not_set",
			prepare: func(request *models.CreateAdvertisementRequest) {
				delete(request.Fields, "price")
			},
		},
		{
			name: "quantity_100",
			prepare: func(request *models.CreateAdvertisementRequest) {
				request.Fields["quantity"] = "100"
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			request := fixtures.AdvertisementRequest(t)
			tc.prepare(&request)

			created := advertisements.CreateAdvertisement(t, user.Token, request)
			assertions.AdvertisementMatchesRequest(t, created, request, user.User.Id)

			actual := advertisements.GetAdvertisement(t, user.Token, created.ID)
			assertions.AdvertisementMatchesRequest(t, actual, request, user.User.Id)
		})
	}

	t.Run("one_photo", func(t *testing.T) {
		request := fixtures.AdvertisementRequest(t)
		created := advertisements.CreateAdvertisement(t, user.Token, request)

		require.Len(t, advertisements.GetAdvertisementPhotos(t, user.Token, created.ID), 1)
	})

	t.Run("three_photos", func(t *testing.T) {
		request := fixtures.AdvertisementRequest(t)
		request.Photos = append(request.Photos,
			fixtures.PNGPhoto(t, "photo-2.png", color.RGBA{R: 200, G: 90, B: 70, A: 255}),
			fixtures.PNGPhoto(t, "photo-3.png", color.RGBA{R: 70, G: 90, B: 200, A: 255}),
		)

		created := advertisements.CreateAdvertisement(t, user.Token, request)
		require.Len(t, advertisements.GetAdvertisementPhotos(t, user.Token, created.ID), 3)
	})
}
