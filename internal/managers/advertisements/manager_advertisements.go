package advertisements

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"

	client "api-tests-template/internal/client/http/advertisements"
	"api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/utils"
)

// GetAdvertisementsList возвращает объявления, найденные по строке поиска.
func GetAdvertisementsList(t *testing.T, token, text string) models.AdvertisementsList {
	t.Helper()
	return utils.DecodeResponse[models.AdvertisementsList](t,
		client.HttpGetAdvertisements(t, token, text), http.StatusOK)
}

// CreateAdvertisement создаёт объявление и регистрирует cleanup на его удаление.
func CreateAdvertisement(t *testing.T, token string, request models.CreateAdvertisementRequest) models.Advertisement {
	t.Helper()

	response := client.HttpPostAdvertisement(t, token, request)
	body := utils.ReadResponseBody(t, response)
	require.Contains(t, []int{http.StatusOK, http.StatusCreated}, response.StatusCode,
		"Создание объявления должно вернуть 200 или 201; response: %s", body)

	registerAdvertisementCleanup(t, token, body)

	var result models.Advertisement
	require.NoError(t, json.Unmarshal(body, &result), "Некорректный advertisement response: %s", body)
	require.NotEmpty(t, result.ID, "Созданное объявление не содержит ID")
	return result
}

// CreateAdvertisementExpectError проверяет ожидаемую ошибку создания объявления.
func CreateAdvertisementExpectError(
	t *testing.T,
	token string,
	request models.CreateAdvertisementRequest,
	expectedStatus int,
) models.ErrorResponse {
	t.Helper()

	response := client.HttpPostAdvertisement(t, token, request)
	body := utils.ReadResponseBody(t, response)

	if response.StatusCode == http.StatusOK || response.StatusCode == http.StatusCreated {
		registerAdvertisementCleanup(t, token, body)
	}

	require.Equal(t, expectedStatus, response.StatusCode,
		"Неожиданный status создания объявления; response: %s", body)

	var result models.ErrorResponse
	require.NoError(t, json.Unmarshal(body, &result), "Некорректный error response: %s", body)
	return result
}

// GetAdvertisement возвращает объявление по ID.
func GetAdvertisement(t *testing.T, token, id string) models.Advertisement {
	t.Helper()
	return utils.DecodeResponse[models.Advertisement](t,
		client.HttpGetAdvertisement(t, token, id), http.StatusOK)
}

// DeleteAdvertisement удаляет объявление и ожидает HTTP 204.
func DeleteAdvertisement(t *testing.T, token, id string) {
	t.Helper()

	response := client.HttpDeleteAdvertisement(t, token, id)
	body := utils.ReadResponseBody(t, response)
	require.Equal(t, http.StatusNoContent, response.StatusCode,
		"Не удалось удалить объявление %s; response: %s", id, body)
}

// GetAdvertisementPhotos возвращает метаданные фотографий объявления.
func GetAdvertisementPhotos(t *testing.T, token, id string) []models.Photo {
	t.Helper()
	return utils.DecodeResponse[[]models.Photo](t,
		client.HttpGetAdvertisementPhotos(t, token, id), http.StatusOK)
}

// DeleteAdvertisementIfExists удаляет объявление при cleanup и допускает уже удалённый объект.
func DeleteAdvertisementIfExists(t *testing.T, token, id string) {
	t.Helper()

	response := client.HttpDeleteAdvertisement(t, token, id)
	body := utils.ReadResponseBody(t, response)

	if response.StatusCode == http.StatusForbidden {
		check := client.HttpGetAdvertisement(t, token, id)
		checkBody := utils.ReadResponseBody(t, check)
		require.Equal(t, http.StatusNotFound, check.StatusCode,
			"DELETE не подтвердил удаление %s; DELETE: %s; GET: %s", id, body, checkBody)
		return
	}

	require.Contains(t, []int{http.StatusNoContent, http.StatusNotFound}, response.StatusCode,
		"Не удалось очистить объявление %s; response: %s", id, body)
}

// registerAdvertisementCleanup регистрирует удаление объявления, если успешный ответ содержит ID.
func registerAdvertisementCleanup(t *testing.T, token string, body []byte) {
	t.Helper()

	id := gjson.GetBytes(body, "id").String()
	if id == "" {
		return
	}

	t.Cleanup(func() {
		t.Logf("[Cleanup] Удаляем объявление %s", id)
		DeleteAdvertisementIfExists(t, token, id)
	})
}
