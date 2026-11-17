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

// GetAdvertisementsList возвращает список найденных объявлений.
func GetAdvertisementsList(t *testing.T, token, text string) models.AdvertisementsList {
	t.Helper()

	body := utils.ReadResponseBody(t, client.HttpGetAdvertisements(t, token, text))

	var result models.AdvertisementsList
	require.NoError(t, json.Unmarshal(body, &result), "Ответ API должен соответствовать JSON-модели")
	return result
}

// CreateAdvertisement создаёт объявление и регистрирует его очистку.
func CreateAdvertisement(t *testing.T, token string, request models.CreateAdvertisementRequest) models.Advertisement {
	t.Helper()

	response := client.HttpPostAdvertisement(t, token, request)
	body := utils.ReadResponseBody(t, response)

	require.Equal(t, http.StatusCreated, response.StatusCode, "Создание объявления: response=%s", body)

	registerAdvertisementCleanup(t, token, body)

	var result models.Advertisement
	require.NoError(t, json.Unmarshal(body, &result), "Ответ API должен соответствовать JSON-модели")
	require.NotEmpty(t, result.ID, "Поле result.ID не должно быть пустым")

	return result
}

// CreateAdvertisementExpectError проверяет статус неуспешного создания.
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

	require.Equal(t, expectedStatus, response.StatusCode, "response: %s", body)

	var result models.ErrorResponse
	require.NoError(t, json.Unmarshal(body, &result), "Ответ API должен соответствовать JSON-модели")

	return result
}

// GetAdvertisement проверяет получение объявления.
func GetAdvertisement(t *testing.T, token, id string) models.Advertisement {
	t.Helper()

	response := client.HttpGetAdvertisement(t, token, id)
	body := utils.ReadResponseBody(t, response)

	require.Equal(t, http.StatusOK, response.StatusCode, "response: %s", body)

	var result models.Advertisement
	require.NoError(t, json.Unmarshal(body, &result), "Ответ API должен соответствовать JSON-модели")

	return result
}

// DeleteAdvertisement проверяет удаление объявления.
func DeleteAdvertisement(t *testing.T, token, id string) {
	t.Helper()

	response := client.HttpDeleteAdvertisement(t, token, id)
	body := utils.ReadResponseBody(t, response)

	require.Equal(t, http.StatusNoContent, response.StatusCode, "response: %s", body)
}

// GetAdvertisementPhotos возвращает фотографии объявления.
func GetAdvertisementPhotos(t *testing.T, token, id string) []models.Photo {
	t.Helper()

	response := client.HttpGetAdvertisementPhotos(t, token, id)
	body := utils.ReadResponseBody(t, response)

	require.Equal(t, http.StatusOK, response.StatusCode, "response: %s", body)

	var result []models.Photo
	require.NoError(t, json.Unmarshal(body, &result), "Ответ API должен соответствовать JSON-модели")

	return result
}

// DeleteAdvertisementIfExists удаляет тестовое объявление, если оно ещё существует.
func DeleteAdvertisementIfExists(t *testing.T, token, id string) {
	t.Helper()

	response := client.HttpGetAdvertisement(t, token, id)
	body := utils.ReadResponseBody(t, response)
	if response.StatusCode == http.StatusNotFound {
		var result models.ErrorResponse
		require.NoError(t, json.Unmarshal(body, &result), "Ответ очистки объявления должен быть JSON")
		require.Equal(t, "Advertisement not found", result.Message, "Неожиданная ошибка очистки объявления")
		return
	}
	require.Equal(t, http.StatusOK, response.StatusCode, "Проверка объявления перед очисткой: response=%s", body)
	DeleteAdvertisement(t, token, id)
}

// registerAdvertisementCleanup регистрирует удаление созданного объявления.
func registerAdvertisementCleanup(t *testing.T, token string, body []byte) {
	t.Helper()

	id := gjson.GetBytes(body, "id").String()
	if id == "" {
		return
	}

	t.Cleanup(func() {
		DeleteAdvertisementIfExists(t, token, id)
	})
}
