package favorites

import (
	"net/http"
	"testing"

	"api-tests-template/internal/managers/advertisements"
	apiModels "api-tests-template/internal/managers/advertisements/models"
	manager "api-tests-template/internal/managers/favorites"
	"api-tests-template/internal/managers/favorites/models"
	"api-tests-template/tests/fixtures"
)

func TestFavoritesWithoutAuthorization(t *testing.T) {
	t.Log("Precondition: подготовка пользователя и тестовых данных")
	owner := fixtures.Login(t)
	advertisement := advertisements.CreateAdvertisement(t, owner.Token, fixtures.AdvertisementRequest(t))

	operations := []struct {
		name string
		call func(*testing.T, string) apiModels.ErrorResponse
	}{
		{
			name: "get",
			call: func(t *testing.T, token string) apiModels.ErrorResponse {
				return manager.GetFavoritesExpectStatus(t, token, http.StatusUnauthorized)
			},
		},
		{
			name: "add",
			call: func(t *testing.T, token string) apiModels.ErrorResponse {
				return manager.AddFavoriteExpectStatus(
					t,
					token,
					models.AddFavoriteRequest{AdvertisementID: advertisement.ID},
					http.StatusUnauthorized,
				)
			},
		},
		{
			name: "remove",
			call: func(t *testing.T, token string) apiModels.ErrorResponse {
				return manager.RemoveFavoriteExpectStatus(
					t,
					token,
					advertisement.ID,
					http.StatusUnauthorized,
				)
			},
		},
	}

	tokens := []struct {
		name    string
		token   string
		message string
	}{
		{name: "missing_token", message: fixtures.MissingTokenMessage},
		{name: "invalid_token", token: fixtures.InvalidToken, message: fixtures.InvalidAuthorizationMessage},
	}

	for _, operation := range operations {
		t.Run(operation.name, func(t *testing.T) {
			t.Log("Step: " + t.Name())
			for _, tokenCase := range tokens {
				t.Run(tokenCase.name, func(t *testing.T) {
					t.Log("Step: " + t.Name())
					response := operation.call(t, tokenCase.token)
					fixtures.AssertError(t, response, fixtures.UnauthorizedError, tokenCase.message)
					t.Log("Result: получены 401 и точная ошибка авторизации")
				})
			}
		})
	}
	if !t.Failed() {
		t.Log("Result: проверки сценария пройдены")
	}
}
