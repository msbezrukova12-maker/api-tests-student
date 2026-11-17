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

// GetAllFavorites получает список избранного авторизованного пользователя.
func GetAllFavorites(t *testing.T, token string) []models.Favorite {
	t.Helper()

	response := utils.DecodeResponse[models.FavoritesList](t,
		client.HttpGetFavorites(t, token), http.StatusOK)
	return response.Items
}

// AddFavorite добавляет объявление в избранное и регистрирует cleanup.
func AddFavorite(t *testing.T, token, advertisementID string) models.Favorite {
	t.Helper()

	response := client.HttpPostFavorite(t, token, models.AddFavoriteRequest{AdvertisementID: advertisementID})
	body := utils.ReadResponseBody(t, response)
	require.GreaterOrEqual(t, response.StatusCode, http.StatusOK,
		"POST /favorites должен завершиться успешно; response: %s", body)
	require.Less(t, response.StatusCode, http.StatusMultipleChoices,
		"POST /favorites должен завершиться успешно; response: %s", body)

	var result models.Favorite
	require.NoError(t, json.Unmarshal(body, &result), "Некорректный favorite JSON: %s", body)
	require.NotEmpty(t, result.FavoriteID, "Favorite не содержит favorite_id")
	require.NotEmpty(t, result.AddedAt, "Favorite не содержит added_at")
	require.Equal(t, advertisementID, result.Advertisement.ID, "Сохранено другое объявление")

	t.Cleanup(func() {
		t.Logf("[Cleanup] Удаляем favorite для объявления %s", advertisementID)
		RemoveFavoriteIfExists(t, token, advertisementID)
	})

	return result
}

// AddFavoriteExpectStatus проверяет HTTP-статус, который прямо задан в условии теста.
func AddFavoriteExpectStatus(t *testing.T, token string, request any, expectedStatus int) apiModels.ErrorResponse {
	t.Helper()

	response := client.HttpPostFavorite(t, token, request)
	return decodeFavoriteError(t, response, expectedStatus)
}

// AddFavoriteExpectClientError проверяет, что бизнес-ошибка POST /favorites является 4xx и содержит сообщение.
func AddFavoriteExpectClientError(t *testing.T, token string, request any) (int, apiModels.ErrorResponse) {
	t.Helper()

	response := client.HttpPostFavorite(t, token, request)
	body := utils.ReadResponseBody(t, response)
	require.GreaterOrEqual(t, response.StatusCode, http.StatusBadRequest,
		"Ожидалась клиентская ошибка; response: %s", body)
	require.Less(t, response.StatusCode, http.StatusInternalServerError,
		"Ожидалась клиентская ошибка; response: %s", body)

	var result apiModels.ErrorResponse
	require.NoError(t, json.Unmarshal(body, &result), "Некорректный favorite error JSON: %s", body)
	require.NotEmpty(t, result.Error, "API должен вернуть тип ошибки")
	require.NotEmpty(t, result.Message, "API должен вернуть сообщение об ошибке")
	return response.StatusCode, result
}

// RemoveFavorite удаляет существующую запись избранного.
func RemoveFavorite(t *testing.T, token, advertisementID string) {
	t.Helper()

	response := client.HttpDeleteFavorite(t, token, advertisementID)
	body := utils.ReadResponseBody(t, response)
	require.GreaterOrEqual(t, response.StatusCode, http.StatusOK,
		"DELETE /favorites должен завершиться успешно; response: %s", body)
	require.Less(t, response.StatusCode, http.StatusMultipleChoices,
		"DELETE /favorites должен завершиться успешно; response: %s", body)
}

// RemoveFavoriteIfExists удаляет favorite во время cleanup и не падает, если запись уже отсутствует.
func RemoveFavoriteIfExists(t *testing.T, token, advertisementID string) {
	t.Helper()

	response := client.HttpDeleteFavorite(t, token, advertisementID)
	body := utils.ReadResponseBody(t, response)
	if response.StatusCode >= http.StatusOK && response.StatusCode < http.StatusMultipleChoices {
		return
	}
	if response.StatusCode >= http.StatusBadRequest && response.StatusCode < http.StatusInternalServerError {
		t.Logf("[Cleanup] Favorite уже отсутствует: status=%d, response=%s", response.StatusCode, body)
		return
	}
	require.Fail(t, "Не удалось очистить favorite",
		"status=%d, response=%s", response.StatusCode, body)
}

// RemoveFavoriteExpectStatus проверяет HTTP-статус, который прямо задан в условии теста.
func RemoveFavoriteExpectStatus(
	t *testing.T,
	token string,
	advertisementID string,
	expectedStatus int,
) apiModels.ErrorResponse {
	t.Helper()

	response := client.HttpDeleteFavorite(t, token, advertisementID)
	return decodeFavoriteError(t, response, expectedStatus)
}

// RemoveFavoriteExpectClientError проверяет, что бизнес-ошибка DELETE /favorites/{id} является 4xx.
func RemoveFavoriteExpectClientError(
	t *testing.T,
	token string,
	advertisementID string,
) (int, apiModels.ErrorResponse) {
	t.Helper()

	response := client.HttpDeleteFavorite(t, token, advertisementID)
	body := utils.ReadResponseBody(t, response)
	require.GreaterOrEqual(t, response.StatusCode, http.StatusBadRequest,
		"Ожидалась клиентская ошибка; response: %s", body)
	require.Less(t, response.StatusCode, http.StatusInternalServerError,
		"Ожидалась клиентская ошибка; response: %s", body)

	var result apiModels.ErrorResponse
	require.NoError(t, json.Unmarshal(body, &result), "Некорректный favorite error JSON: %s", body)
	require.NotEmpty(t, result.Error, "API должен вернуть тип ошибки")
	require.NotEmpty(t, result.Message, "API должен вернуть сообщение об ошибке")
	return response.StatusCode, result
}

// GetFavoritesExpectStatus проверяет HTTP-статус, который прямо задан в условии теста.
func GetFavoritesExpectStatus(t *testing.T, token string, expectedStatus int) apiModels.ErrorResponse {
	t.Helper()

	response := client.HttpGetFavorites(t, token)
	return decodeFavoriteError(t, response, expectedStatus)
}

// decodeFavoriteError проверяет ожидаемый статус и наличие структурированного сообщения об ошибке.
func decodeFavoriteError(t *testing.T, response *http.Response, expectedStatus int) apiModels.ErrorResponse {
	t.Helper()

	body := utils.ReadResponseBody(t, response)
	require.Equal(t, expectedStatus, response.StatusCode,
		"Неожиданный HTTP status; response: %s", body)

	var result apiModels.ErrorResponse
	require.NoError(t, json.Unmarshal(body, &result), "Некорректный error JSON: %s", body)
	require.NotEmpty(t, result.Error, "API должен вернуть тип ошибки")
	require.NotEmpty(t, result.Message, "API должен вернуть сообщение об ошибке")
	return result
}

