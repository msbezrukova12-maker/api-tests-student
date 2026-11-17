package auth

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/client/http/auth"
	"api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/utils"
)

func Login(t *testing.T, username, password string) models.LoginOkResponse {
	t.Helper()

	require.NotEmpty(t, username)
	require.NotEmpty(t, password)

	request := models.LoginRequest{
		Email:    username,
		Password: password,
	}

	response := auth.HttpPostAuthLogin(t, request)
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	require.NoError(t, err)

	var result models.LoginOkResponse
	require.NoError(t, json.Unmarshal(body, &result))

	return result
}

func RegisterOrLogin(t *testing.T, request models.RegisterRequest) models.LoginOkResponse {
	t.Helper()

	response := auth.HttpPostAuthRegister(t, request)
	body := utils.ReadResponseBody(t, response)

	if response.StatusCode >= http.StatusOK && response.StatusCode < http.StatusMultipleChoices {
		var result models.LoginOkResponse
		require.NoError(t, json.Unmarshal(body, &result))
		return result
	}

	return Login(t, request.Email, request.Password)
}
