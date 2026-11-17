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

func GetAllFavorites(t *testing.T, token string) []models.Favorite {
	t.Helper()

	response := utils.DecodeResponse[models.FavoritesList](
		t,
		client.HttpGetFavorites(t, token),
		http.StatusOK,
	)

	return response.Items
}

func AddFavorite(t *testing.T, token, advertisementID string) models.Favorite {
	t.Helper()

	response := client.HttpPostFavorite(
		t,
		token,
		models.AddFavoriteRequest{AdvertisementID: advertisementID},
	)
	body := utils.ReadResponseBody(t, response)

	require.GreaterOrEqual(t, response.StatusCode, http.StatusOK, "response: %s", body)
	require.Less(t, response.StatusCode, http.StatusMultipleChoices, "response: %s", body)

	var result models.Favorite
	require.NoError(t, json.Unmarshal(body, &result))
	require.Equal(t, advertisementID, result.Advertisement.ID)

	t.Cleanup(func() {
		RemoveFavoriteIfExists(t, token, advertisementID)
	})

	return result
}

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

func AddFavoriteExpectClientError(
	t *testing.T,
	token string,
	request any,
) (int, apiModels.ErrorResponse) {
	t.Helper()

	response := client.HttpPostFavorite(t, token, request)
	body := utils.ReadResponseBody(t, response)

	require.GreaterOrEqual(t, response.StatusCode, http.StatusBadRequest, "response: %s", body)
	require.Less(t, response.StatusCode, http.StatusInternalServerError, "response: %s", body)

	var result apiModels.ErrorResponse
	require.NoError(t, json.Unmarshal(body, &result))
	require.NotEmpty(t, result.Error)
	require.NotEmpty(t, result.Message)

	return response.StatusCode, result
}

func RemoveFavorite(t *testing.T, token, advertisementID string) {
	t.Helper()

	response := client.HttpDeleteFavorite(t, token, advertisementID)
	body := utils.ReadResponseBody(t, response)

	require.GreaterOrEqual(t, response.StatusCode, http.StatusOK, "response: %s", body)
	require.Less(t, response.StatusCode, http.StatusMultipleChoices, "response: %s", body)
}

func RemoveFavoriteIfExists(t *testing.T, token, advertisementID string) {
	t.Helper()

	response := client.HttpDeleteFavorite(t, token, advertisementID)
	body := utils.ReadResponseBody(t, response)

	if response.StatusCode >= http.StatusOK && response.StatusCode < http.StatusInternalServerError {
		return
	}

	require.Fail(t, "cleanup failed", "status=%d, response=%s", response.StatusCode, body)
}

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

func RemoveFavoriteExpectClientError(
	t *testing.T,
	token string,
	advertisementID string,
) (int, apiModels.ErrorResponse) {
	t.Helper()

	response := client.HttpDeleteFavorite(t, token, advertisementID)
	body := utils.ReadResponseBody(t, response)

	require.GreaterOrEqual(t, response.StatusCode, http.StatusBadRequest, "response: %s", body)
	require.Less(t, response.StatusCode, http.StatusInternalServerError, "response: %s", body)

	var result apiModels.ErrorResponse
	require.NoError(t, json.Unmarshal(body, &result))
	require.NotEmpty(t, result.Error)
	require.NotEmpty(t, result.Message)

	return response.StatusCode, result
}

func GetFavoritesExpectStatus(t *testing.T, token string, expectedStatus int) apiModels.ErrorResponse {
	t.Helper()

	return decodeFavoriteError(
		t,
		client.HttpGetFavorites(t, token),
		expectedStatus,
	)
}

func decodeFavoriteError(
	t *testing.T,
	response *http.Response,
	expectedStatus int,
) apiModels.ErrorResponse {
	t.Helper()

	body := utils.ReadResponseBody(t, response)
	require.Equal(t, expectedStatus, response.StatusCode, "response: %s", body)

	var result apiModels.ErrorResponse
	require.NoError(t, json.Unmarshal(body, &result))
	require.NotEmpty(t, result.Error)
	require.NotEmpty(t, result.Message)

	return result
}
