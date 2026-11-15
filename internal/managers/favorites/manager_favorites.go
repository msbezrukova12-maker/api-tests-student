package favorites

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	client "api-tests-template/internal/client/http/favorites"
	apiModels "api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/managers/favorites/models"
	"api-tests-template/internal/utils"
)

const favoritesPageSize = 20

// GetAllFavorites получает все страницы избранного и проверяет корректность пагинации.
func GetAllFavorites(t *testing.T, token string) []models.Favorite {
	t.Helper()

	var items []models.Favorite
	for page := 1; ; page++ {
		response := utils.DecodeResponse[models.FavoritesList](t,
			client.HttpGetFavorites(t, token, page, favoritesPageSize), http.StatusOK)

		require.Equal(t, page, response.Page, "API вернул неправильный номер страницы")
		require.Equal(t, favoritesPageSize, response.PageSize, "API вернул неправильный page_size")
		items = append(items, response.Items...)

		if page >= response.TotalPages {
			require.Len(t, items, response.Total, "Пагинация потеряла или продублировала элементы")
			return items
		}

		require.NotEmpty(t, response.Items, "Получена пустая страница до последней")
	}
}

// AddFavorite добавляет объявление в избранное и регистрирует cleanup.
func AddFavorite(t *testing.T, token, advertisementID string) models.Favorite {
	t.Helper()

	body := postFavorite(t, token, models.AddFavoriteRequest{AdvertisementID: advertisementID}, http.StatusCreated)

	var result models.Favorite
	require.NoError(t, json.Unmarshal(body, &result), "Некорректный favorite JSON: %s", body)
	require.NotEmpty(t, result.FavoriteID, "Favorite не содержит favorite_id")
	require.NotEmpty(t, result.AddedAt, "Favorite не содержит added_at")
	require.Equal(t, advertisementID, result.Advertisement.ID, "Сохранено другое объявление")
	return result
}

// AddFavoriteExpectError проверяет ожидаемую ошибку POST /favorites.
func AddFavoriteExpectError(t *testing.T, token string, request any, expectedStatus int) apiModels.ErrorResponse {
	t.Helper()

	body := postFavorite(t, token, request, expectedStatus)

	var result apiModels.ErrorResponse
	require.NoError(t, json.Unmarshal(body, &result), "Некорректный favorite error JSON: %s", body)
	return result
}

// postFavorite отправляет POST /favorites и регистрирует cleanup при неожиданном успешном создании.
func postFavorite(t *testing.T, token string, request any, expectedStatus int) []byte {
	t.Helper()

	response := client.HttpPostFavorite(t, token, request)
	body := utils.ReadResponseBody(t, response)

	if response.StatusCode == http.StatusCreated {
		var identity struct {
			Advertisement struct {
				ID string `json:"id"`
			} `json:"advertisement"`
		}
		require.NoError(t, json.Unmarshal(body, &identity), "Некорректный favorite response: %s", body)

		if identity.Advertisement.ID != "" {
			t.Cleanup(func() {
				t.Logf("[Cleanup] Удаляем favorite для объявления %s", identity.Advertisement.ID)
				RemoveFavoriteIfExists(t, token, identity.Advertisement.ID)
			})
		}
	}

	require.Equal(t, expectedStatus, response.StatusCode,
		"Неожиданный POST /favorites status; response: %s", body)
	return body
}

// RemoveFavorite удаляет существующую запись избранного.
func RemoveFavorite(t *testing.T, token, advertisementID string) {
	t.Helper()

	response := client.HttpDeleteFavorite(t, token, advertisementID)
	body := utils.ReadResponseBody(t, response)
	require.Equal(t, http.StatusNoContent, response.StatusCode,
		"Не удалось удалить favorite %s; response: %s", advertisementID, body)
}

// RemoveFavoriteIfExists удаляет favorite во время cleanup и допускает его отсутствие.
func RemoveFavoriteIfExists(t *testing.T, token, advertisementID string) {
	t.Helper()

	response := client.HttpDeleteFavorite(t, token, advertisementID)
	body := utils.ReadResponseBody(t, response)
	require.Contains(t, []int{http.StatusNoContent, http.StatusNotFound}, response.StatusCode,
		"Не удалось очистить favorite %s; response: %s", advertisementID, body)
}

// RemoveFavoriteExpectError проверяет ожидаемую ошибку DELETE /favorites/{id}.
func RemoveFavoriteExpectError(
	t *testing.T,
	token string,
	advertisementID string,
	expectedStatus int,
) apiModels.ErrorResponse {
	t.Helper()
	return utils.DecodeResponse[apiModels.ErrorResponse](t,
		client.HttpDeleteFavorite(t, token, advertisementID), expectedStatus)
}

// GetFavoritesExpectError проверяет ожидаемую ошибку GET /favorites.
func GetFavoritesExpectError(t *testing.T, token string, expectedStatus int) apiModels.ErrorResponse {
	t.Helper()
	return utils.DecodeResponse[apiModels.ErrorResponse](t,
		client.HttpGetFavorites(t, token, 1, favoritesPageSize), expectedStatus)
}
