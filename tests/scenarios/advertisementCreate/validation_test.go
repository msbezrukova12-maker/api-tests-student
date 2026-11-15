package advertisementCreate

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/tests/fixtures"
)

const (
	maximumTitleLength     = 50
	invalidRequestError    = "invalid_request"
	priceValidationMessage = "price must be between 0 and 1000000 inclusive"
)

// TestCreateAdvertisementInvalidParameters проверяет обязательные поля и основные ошибки валидации.
func TestCreateAdvertisementInvalidParameters(t *testing.T) {
	user := fixtures.Login(t)

	cases := []struct {
		name        string
		prepare     func(map[string]string)
		removePhoto bool
		message     string
	}{
		{
			name: "missing_title",
			prepare: func(fields map[string]string) {
				delete(fields, "title")
			},
			message: "title is required",
		},
		{
			name: "empty_title",
			prepare: func(fields map[string]string) {
				fields["title"] = ""
			},
			message: "title is required",
		},
		{
			name: "title_above_limit",
			prepare: func(fields map[string]string) {
				fields["title"] = strings.Repeat("a", maximumTitleLength+1)
			},
			message: fmt.Sprintf(
				"title must not exceed %d characters provided %d",
				maximumTitleLength,
				maximumTitleLength+1,
			),
		},
		{
			name: "negative_price",
			prepare: func(fields map[string]string) {
				fields["price"] = "-1"
			},
			message: priceValidationMessage,
		},
		{
			name:        "missing_photos",
			prepare:     func(map[string]string) {},
			removePhoto: true,
			message:     "at least 1 photo is required",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Logf("[Precondition] Подготавливаем невалидный request: %s", tc.name)
			request := fixtures.AdvertisementRequest(t)
			tc.prepare(request.Fields)
			if tc.removePhoto {
				request.Photos = nil
			}

			t.Log("[Step] Отправляем невалидный POST /advertisement")
			response := advertisements.CreateAdvertisementExpectError(
				t,
				user.Token,
				request,
				http.StatusBadRequest,
			)

			require.Equal(t, invalidRequestError, response.Error, "Неправильный тип validation error")
			require.Equal(t, tc.message, response.Message, "Неправильное validation message")
			t.Logf("[Result] Получен HTTP 400: %s / %s", response.Error, response.Message)
		})
	}
}
