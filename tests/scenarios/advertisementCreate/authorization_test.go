package advertisementCreate

import (
	"net/http"
	"testing"

	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/tests/fixtures"
)

func TestCreateAdvertisementWithoutAuthorization(t *testing.T) {
	t.Log("Precondition: подготовка пользователя и тестовых данных")
	cases := []struct {
		name    string
		token   string
		message string
	}{
		{name: "missing_token", message: fixtures.MissingTokenMessage},
		{name: "invalid_token", token: fixtures.InvalidToken, message: fixtures.InvalidAuthorizationMessage},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Log("Step: " + t.Name())
			fixtures.LoadEnvironment(t)
			request := fixtures.AdvertisementRequest(t)

			response := advertisements.CreateAdvertisementExpectError(
				t,
				tc.token,
				request,
				http.StatusUnauthorized,
			)

			fixtures.AssertError(t, response, fixtures.UnauthorizedError, tc.message)
			t.Log("Result: получены 401 и точная ошибка авторизации")
		})
	}
	if !t.Failed() {
		t.Log("Result: проверки сценария пройдены")
	}
}
