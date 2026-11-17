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

const invalidToken = "incorrect_token"

// TestFavoritesWithoutAuthorization проверяет 401 для ручек избранного без валидной авторизации.
func TestFavoritesWithoutAuthorization(t *testing.T) {
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
		name  string
		token string
	}{
		{name: "missing_token"},
		{name: "invalid_token", token: invalidToken},
	}

	for _, operation := range operations {
		t.Run(operation.name, func(t *testing.T) {
			for _, tokenCase := range tokens {
				t.Run(tokenCase.name, func(t *testing.T) {
					t.Logf("[Step] Вызываем %s без валидной авторизации: %s", operation.name, tokenCase.name)
					response := operation.call(t, tokenCase.token)
					t.Logf("[Result] HTTP 401, error=%q, message=%q", response.Error, response.Message)
				})
			}
		})
	}
}
