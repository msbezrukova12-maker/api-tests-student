package advertisementCreate

import (
	"net/http"
	"strings"
	"testing"

	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/tests/fixtures"
)

func TestCreateAdvertisementInvalidParameters(t *testing.T) {
	t.Log("Precondition: подготовка пользователя и тестовых данных")
	user := fixtures.Login(t)

	cases := []struct {
		name    string
		message string
		prepare func(*models.CreateAdvertisementRequest)
	}{
		{
			name:    "missing_title",
			message: fixtures.RequiredTitleMessage,
			prepare: func(request *models.CreateAdvertisementRequest) {
				delete(request.Fields, "title")
			},
		},
		{
			name:    "empty_title",
			message: fixtures.RequiredTitleMessage,
			prepare: func(request *models.CreateAdvertisementRequest) {
				request.Fields["title"] = ""
			},
		},
		{
			name:    "title_51_characters",
			message: "title must not exceed 50 characters provided 51",
			prepare: func(request *models.CreateAdvertisementRequest) {
				request.Fields["title"] = strings.Repeat("a", 51)
			},
		},
		{
			name:    "missing_description",
			message: "description is required",
			prepare: func(request *models.CreateAdvertisementRequest) {
				delete(request.Fields, "description")
			},
		},
		{
			name:    "description_501_characters",
			message: "description must be between 1 and 500 characters",
			prepare: func(request *models.CreateAdvertisementRequest) {
				request.Fields["description"] = strings.Repeat("b", 501)
			},
		},
		{
			name:    "negative_price",
			message: fixtures.PriceValidationMessage,
			prepare: func(request *models.CreateAdvertisementRequest) {
				request.Fields["price"] = "-1"
			},
		},
		{
			name:    "price_above_1000000",
			message: fixtures.PriceValidationMessage,
			prepare: func(request *models.CreateAdvertisementRequest) {
				request.Fields["price"] = "1000001"
			},
		},
		{
			name:    "quantity_101",
			message: "quantity must be between 1 and 100 inclusive",
			prepare: func(request *models.CreateAdvertisementRequest) {
				request.Fields["quantity"] = "101"
			},
		},
		{
			name:    "missing_photos",
			message: "at least 1 photo is required",
			prepare: func(request *models.CreateAdvertisementRequest) {
				request.Photos = nil
			},
		},
		{
			name:    "unsupported_photo_format",
			message: `photo 1: file type not allowed: extension ".txt"`,
			prepare: func(request *models.CreateAdvertisementRequest) {
				request.Photos = []models.PhotoUpload{
					{Name: "photo.txt", Content: []byte("This is a text file.")},
				}
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Log("Step: " + t.Name())
			request := fixtures.AdvertisementRequest(t)
			tc.prepare(&request)

			response := advertisements.CreateAdvertisementExpectError(
				t,
				user.Token,
				request,
				http.StatusBadRequest,
			)

			fixtures.AssertError(t, response, fixtures.InvalidRequestError, tc.message)
			t.Log("Result: status и ошибка соответствуют ожиданиям")
		})
	}
	if !t.Failed() {
		t.Log("Result: проверки сценария пройдены")
	}
}
