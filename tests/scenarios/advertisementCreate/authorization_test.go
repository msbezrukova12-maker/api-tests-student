package advertisementCreate

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/tests/fixtures"
)

const invalidToken = "incorrect_token"

// TestCreateAdvertisementWithoutAuthorization проверяет 401 для отсутствующего и неверного токена.
func TestCreateAdvertisementWithoutAuthorization(t *testing.T) {
	cases := []struct {
		name  string
		token string
	}{
		{name: "missing_token"},
		{name: "invalid_token", token: invalidToken},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fixtures.LoadEnvironment(t)
			request := fixtures.AdvertisementRequest(t)

			t.Logf("[Step] Создаём объявление без валидной авторизации: %s", tc.name)
			response := advertisements.CreateAdvertisementExpectError(
				t,
				tc.token,
				request,
				http.StatusUnauthorized,
			)

			require.NotEmpty(t, response.Error, "API должен вернуть тип ошибки авторизации")
			require.NotEmpty(t, response.Message, "API должен вернуть сообщение об ошибке авторизации")
			t.Logf("[Result] HTTP 401, error=%q, message=%q", response.Error, response.Message)
		})
	}
}
