package advertisementCreate

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/tests/fixtures"
)

const invalidToken = "incorrect_token"

// TestCreateAdvertisementWithoutAuthorization проверяет отсутствие и невалидность bearer-токена.
func TestCreateAdvertisementWithoutAuthorization(t *testing.T) {
	cases := []struct {
		name    string
		token   string
		message string
	}{
		{name: "missing_token", message: "Authorization header required"},
		{name: "invalid_token", token: invalidToken, message: "Invalid or expired token"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fixtures.LoadEnvironment(t)
			request := fixtures.AdvertisementRequest(t)

			t.Logf("[Step] Создаём объявление с auth case=%s", tc.name)
			response := advertisements.CreateAdvertisementExpectError(
				t,
				tc.token,
				request,
				http.StatusUnauthorized,
			)

			require.Equal(t, "unauthorized", response.Error, "Неправильный auth error")
			require.Equal(t, tc.message, response.Message, "Неправильное auth message")
			t.Log("[Result] Получен ожидаемый HTTP 401")
		})
	}
}
