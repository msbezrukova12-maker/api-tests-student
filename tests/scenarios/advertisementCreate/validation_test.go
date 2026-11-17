package advertisementCreate

import (
	"image/color"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/tests/fixtures"
)

// TestCreateAdvertisementInvalidParameters автоматизирует негативные кейсы из ЛР2 и текущего задания Avito.
func TestCreateAdvertisementInvalidParameters(t *testing.T) {
	user := fixtures.Login(t)

	cases := []struct {
		name    string
		prepare func(*testing.T, *models.CreateAdvertisementRequest)
	}{
		{
			name: "missing_title",
			prepare: func(_ *testing.T, request *models.CreateAdvertisementRequest) {
				delete(request.Fields, "title")
			},
		},
		{
			name: "empty_title",
			prepare: func(_ *testing.T, request *models.CreateAdvertisementRequest) {
				request.Fields["title"] = ""
			},
		},
		{
			name: "title_51_characters",
			prepare: func(_ *testing.T, request *models.CreateAdvertisementRequest) {
				request.Fields["title"] = strings.Repeat("a", 51)
			},
		},
		{
			name: "missing_description",
			prepare: func(_ *testing.T, request *models.CreateAdvertisementRequest) {
				delete(request.Fields, "description")
			},
		},
		{
			name: "description_only_spaces",
			prepare: func(_ *testing.T, request *models.CreateAdvertisementRequest) {
				request.Fields["description"] = "   "
			},
		},
		{
			name: "description_501_characters",
			prepare: func(_ *testing.T, request *models.CreateAdvertisementRequest) {
				request.Fields["description"] = strings.Repeat("b", 501)
			},
		},
		{
			name: "negative_price",
			prepare: func(_ *testing.T, request *models.CreateAdvertisementRequest) {
				request.Fields["price"] = "-1"
			},
		},
		{
			name: "price_above_1000000",
			prepare: func(_ *testing.T, request *models.CreateAdvertisementRequest) {
				request.Fields["price"] = "1000001"
			},
		},
		{
			name: "quantity_101",
			prepare: func(_ *testing.T, request *models.CreateAdvertisementRequest) {
				request.Fields["quantity"] = "101"
			},
		},
		{
			name: "missing_photos",
			prepare: func(_ *testing.T, request *models.CreateAdvertisementRequest) {
				request.Photos = nil
			},
		},
		{
			name: "four_photos",
			prepare: func(t *testing.T, request *models.CreateAdvertisementRequest) {
				request.Photos = append(request.Photos,
					fixtures.PNGPhoto(t, "photo-2.png", color.RGBA{R: 200, G: 80, B: 80, A: 255}),
					fixtures.PNGPhoto(t, "photo-3.png", color.RGBA{R: 80, G: 200, B: 80, A: 255}),
					fixtures.PNGPhoto(t, "photo-4.png", color.RGBA{R: 80, G: 80, B: 200, A: 255}),
				)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			request := fixtures.AdvertisementRequest(t)
			tc.prepare(t, &request)

			t.Logf("[Step] Отправляем невалидный POST /advertisement: %s", tc.name)
			response := advertisements.CreateAdvertisementExpectError(
				t,
				user.Token,
				request,
				http.StatusBadRequest,
			)

			require.NotEmpty(t, response.Error, "API должен вернуть тип ошибки")
			require.NotEmpty(t, response.Message, "API должен вернуть сообщение об ошибке")
			t.Logf("[Result] HTTP 400, error=%q, message=%q", response.Error, response.Message)
		})
	}
}
