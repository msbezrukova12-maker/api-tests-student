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

func GetAdvertisementsList(t *testing.T, token, text string) models.AdvertisementsList {
	t.Helper()

	body := utils.ReadResponseBody(t, client.HttpGetAdvertisements(t, token, text))

	var result models.AdvertisementsList
	require.NoError(t, json.Unmarshal(body, &result))
	return result
}

func CreateAdvertisement(t *testing.T, token string, request models.CreateAdvertisementRequest) models.Advertisement {
	t.Helper()

	response := client.HttpPostAdvertisement(t, token, request)
	body := utils.ReadResponseBody(t, response)

	require.Contains(t, []int{http.StatusOK, http.StatusCreated}, response.StatusCode,
		"response: %s", body)

	registerAdvertisementCleanup(t, token, body)

	var result models.Advertisement
	require.NoError(t, json.Unmarshal(body, &result))
	require.NotEmpty(t, result.ID)

	return result
}

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
	require.NoError(t, json.Unmarshal(body, &result))

	return result
}

func GetAdvertisement(t *testing.T, token, id string) models.Advertisement {
	t.Helper()

	return utils.DecodeResponse[models.Advertisement](
		t,
		client.HttpGetAdvertisement(t, token, id),
		http.StatusOK,
	)
}

func DeleteAdvertisement(t *testing.T, token, id string) {
	t.Helper()

	response := client.HttpDeleteAdvertisement(t, token, id)
	body := utils.ReadResponseBody(t, response)

	require.Equal(t, http.StatusNoContent, response.StatusCode, "response: %s", body)
}

func GetAdvertisementPhotos(t *testing.T, token, id string) []models.Photo {
	t.Helper()

	return utils.DecodeResponse[[]models.Photo](
		t,
		client.HttpGetAdvertisementPhotos(t, token, id),
		http.StatusOK,
	)
}

func DeleteAdvertisementIfExists(t *testing.T, token, id string) {
	t.Helper()

	response := client.HttpDeleteAdvertisement(t, token, id)
	body := utils.ReadResponseBody(t, response)

	if response.StatusCode >= http.StatusNoContent && response.StatusCode < http.StatusInternalServerError {
		return
	}

	require.Fail(t, "cleanup failed", "status=%d, response=%s", response.StatusCode, body)
}

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
