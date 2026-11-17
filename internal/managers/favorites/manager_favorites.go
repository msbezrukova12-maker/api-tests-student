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

// GetAllFavorites возвращает первую страницу избранного тестового пользователя.
func GetAllFavorites(t *testing.T, token string) []models.Favorite {
	t.Helper()

	response := client.HttpGetFavorites(t, token)
	body := utils.ReadResponseBody(t, response)

	require.Equal(t, http.StatusOK, response.StatusCode, "response: %s", body)

	var result models.FavoritesList
	require.NoError(t, json.Unmarshal(body, &result), "Ответ API должен соответствовать JSON-модели")

	return result.Items
}

// AddFavorite добавляет объявление в избранное и регистрирует очистку.
func AddFavorite(t *testing.T, token, advertisementID string) models.Favorite {
	t.Helper()

	response := client.HttpPostFavorite(
		t,
		token,
		models.AddFavoriteRequest{AdvertisementID: advertisementID},
	)
	body := utils.ReadResponseBody(t, response)

	require.Equal(t, http.StatusCreated, response.StatusCode, "Добавление избранного: response=%s", body)

	var result models.Favorite
	require.NoError(t, json.Unmarshal(body, &result), "Ответ API должен соответствовать JSON-модели")
	require.Equal(t, advertisementID, result.Advertisement.ID, "Значение result.Advertisement.ID не совпадает с ожидаемым")

	t.Cleanup(func() {
		RemoveFavoriteIfExists(t, token, advertisementID)
	})

	return result
}

// AddFavoriteExpectStatus проверяет статус ошибки добавления в избранное.
func AddFavoriteExpectStatus(
	t *testing.T,
	token string,
	request any,
	expectedStatus int,
) apiModels.ErrorResponse {
	t.Helper()

	return decodeFavoriteError(
		t,
		client.HttpPostFavorite(t, token, request),
		expectedStatus,
	)
}

// RemoveFavorite проверяет удаление из избранного.
func RemoveFavorite(t *testing.T, token, advertisementID string) {
	t.Helper()

	response := client.HttpDeleteFavorite(t, token, advertisementID)
	body := utils.ReadResponseBody(t, response)

	require.Equal(t, http.StatusNoContent, response.StatusCode, "Удаление избранного: response=%s", body)
}

// RemoveFavoriteIfExists очищает избранное, допуская уже удалённую запись.
func RemoveFavoriteIfExists(t *testing.T, token, advertisementID string) {
	t.Helper()

	response := client.HttpDeleteFavorite(t, token, advertisementID)
	body := utils.ReadResponseBody(t, response)

	if response.StatusCode == http.StatusNoContent {
		return
	}

	require.Equal(t, http.StatusNotFound, response.StatusCode, "Очистка избранного: response=%s", body)
	var result apiModels.ErrorResponse
	require.NoError(t, json.Unmarshal(body, &result), "Ответ очистки избранного должен быть JSON")
	require.Equal(t, "Favorite entry not found", result.Message, "Неожиданная ошибка очистки избранного")
}

// RemoveFavoriteExpectStatus проверяет статус ошибки удаления из избранного.
func RemoveFavoriteExpectStatus(
	t *testing.T,
	token string,
	advertisementID string,
	expectedStatus int,
) apiModels.ErrorResponse {
	t.Helper()

	return decodeFavoriteError(
		t,
		client.HttpDeleteFavorite(t, token, advertisementID),
		expectedStatus,
	)
}

// GetFavoritesExpectStatus проверяет статус ошибки получения избранного.
func GetFavoritesExpectStatus(t *testing.T, token string, expectedStatus int) apiModels.ErrorResponse {
	t.Helper()

	return decodeFavoriteError(
		t,
		client.HttpGetFavorites(t, token),
		expectedStatus,
	)
}

// decodeFavoriteError проверяет статус и разбирает ошибочный ответ.
func decodeFavoriteError(
	t *testing.T,
	response *http.Response,
	expectedStatus int,
) apiModels.ErrorResponse {
	t.Helper()

	body := utils.ReadResponseBody(t, response)
	require.Equal(t, expectedStatus, response.StatusCode, "response: %s", body)

	var result apiModels.ErrorResponse
	require.NoError(t, json.Unmarshal(body, &result), "Ответ API должен соответствовать JSON-модели")
	require.NotEmpty(t, result.Error, "Поле result.Error не должно быть пустым")
	require.NotEmpty(t, result.Message, "Поле result.Message не должно быть пустым")

	return result
}
