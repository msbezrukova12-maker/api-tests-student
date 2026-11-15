package favorites

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/advertisements"
	apiModels "api-tests-template/internal/managers/advertisements/models"
	manager "api-tests-template/internal/managers/favorites"
	"api-tests-template/internal/managers/favorites/models"
	"api-tests-template/tests/fixtures"
)

const invalidToken = "incorrect_token"

// TestFavoritesWithoutAuthorization проверяет GET, POST и DELETE favorites без валидной авторизации.
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
				return manager.GetFavoritesExpectError(t, token, http.StatusUnauthorized)
			},
		},
		{
			name: "add",
			call: func(t *testing.T, token string) apiModels.ErrorResponse {
				return manager.AddFavoriteExpectError(
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
				return manager.RemoveFavoriteExpectError(
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
		{name: "missing_token", message: "Authorization header required"},
		{name: "invalid_token", token: invalidToken, message: "Invalid or expired token"},
	}

	for _, operation := range operations {
		t.Run(operation.name, func(t *testing.T) {
			for _, tokenCase := range tokens {
				t.Run(tokenCase.name, func(t *testing.T) {
					t.Logf("[Step] Вызываем %s favorites с %s", operation.name, tokenCase.name)
					response := operation.call(t, tokenCase.token)
					require.Equal(t, "unauthorized", response.Error, "Неправильный auth error")
					require.Equal(t, tokenCase.message, response.Message, "Неправильное auth message")
				})
			}
		})
	}
}
