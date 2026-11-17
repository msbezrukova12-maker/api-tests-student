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

// TestCreateAdvertisementValidBoundaries проверяет валидные граничные значения из ЛР2 и ТЗ Avito.
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

			t.Logf("[Step] Создаём объявление с валидной границей %s", tc.name)
			created := advertisements.CreateAdvertisement(t, user.Token, request)
			assertions.AdvertisementMatchesRequest(t, created, request, user.User.Id)

			actual := advertisements.GetAdvertisement(t, user.Token, created.ID)
			assertions.AdvertisementMatchesRequest(t, actual, request, user.User.Id)

			if tc.name == "price_not_set" {
				require.Zero(t, created.Price, "Если цена не передана, backend должен сохранить нулевую цену")
				require.Zero(t, actual.Price, "GET должен вернуть нулевую цену для объявления без цены")
			}

			t.Log("[Result] Валидная граница принята и сохранена")
		})
	}

	t.Run("one_photo", func(t *testing.T) {
		request := fixtures.AdvertisementRequest(t)
		require.Len(t, request.Photos, 1)

		created := advertisements.CreateAdvertisement(t, user.Token, request)
		require.Len(t, advertisements.GetAdvertisementPhotos(t, user.Token, created.ID), 1,
			"Объявление с одной фотографией должно содержать одну фотографию")
	})

	t.Run("three_photos", func(t *testing.T) {
		request := fixtures.AdvertisementRequest(t)
		request.Photos = append(request.Photos,
			fixtures.PNGPhoto(t, "photo-2.png", color.RGBA{R: 200, G: 90, B: 70, A: 255}),
			fixtures.PNGPhoto(t, "photo-3.png", color.RGBA{R: 70, G: 90, B: 200, A: 255}),
		)

		created := advertisements.CreateAdvertisement(t, user.Token, request)
		require.Len(t, advertisements.GetAdvertisementPhotos(t, user.Token, created.ID), 3,
			"Объявление с тремя фотографиями должно содержать три фотографии")
	})
}
