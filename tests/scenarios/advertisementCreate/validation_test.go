package advertisementCreate

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/tests/fixtures"
)

func TestCreateAdvertisementInvalidParameters(t *testing.T) {
	user := fixtures.Login(t)

	cases := []struct {
		name    string
		prepare func(*models.CreateAdvertisementRequest)
	}{
		{
			name: "missing_title",
			prepare: func(request *models.CreateAdvertisementRequest) {
				delete(request.Fields, "title")
			},
		},
		{
			name: "empty_title",
			prepare: func(request *models.CreateAdvertisementRequest) {
				request.Fields["title"] = ""
			},
		},
		{
			name: "title_51_characters",
			prepare: func(request *models.CreateAdvertisementRequest) {
				request.Fields["title"] = strings.Repeat("a", 51)
			},
		},
		{
			name: "missing_description",
			prepare: func(request *models.CreateAdvertisementRequest) {
				delete(request.Fields, "description")
			},
		},
		{
			name: "description_501_characters",
			prepare: func(request *models.CreateAdvertisementRequest) {
				request.Fields["description"] = strings.Repeat("b", 501)
			},
		},
		{
			name: "negative_price",
			prepare: func(request *models.CreateAdvertisementRequest) {
				request.Fields["price"] = "-1"
			},
		},
		{
			name: "price_above_1000000",
			prepare: func(request *models.CreateAdvertisementRequest) {
				request.Fields["price"] = "1000001"
			},
		},
		{
			name: "quantity_101",
			prepare: func(request *models.CreateAdvertisementRequest) {
				request.Fields["quantity"] = "101"
			},
		},
		{
			name: "missing_photos",
			prepare: func(request *models.CreateAdvertisementRequest) {
				request.Photos = nil
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			request := fixtures.AdvertisementRequest(t)
			tc.prepare(&request)

			response := advertisements.CreateAdvertisementExpectError(
				t,
				user.Token,
				request,
				http.StatusBadRequest,
			)

			require.NotEmpty(t, response.Error)
			require.NotEmpty(t, response.Message)
		})
	}
}
