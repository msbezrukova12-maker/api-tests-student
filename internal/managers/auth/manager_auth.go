package auth

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/client/http/auth"
	"api-tests-template/internal/managers/auth/models"
)

// Login выполняет авторизацию пользователя и возвращает токен с данными пользователя, ожидает 200
func Login(t *testing.T, username, password string) models.LoginOkResponse {
	require.NotEmpty(t, username, "Отсутствует логин пользователя")
	require.NotEmpty(t, password, "Отсутствует пароль пользователя")

	request := models.LoginRequest{
		Email:    username,
		Password: password,
	}

	loginRawResponse := auth.HttpPostAuthLogin(t, request)
	defer loginRawResponse.Body.Close()

	respArray, errResponse := io.ReadAll(loginRawResponse.Body)
	require.NoError(t, errResponse)

	loginResponse := models.LoginOkResponse{}
	err := json.Unmarshal(respArray, &loginResponse)
	require.NoError(t, err, "Операция должна завершиться без ошибки")

	return loginResponse
}

// Register регистрирует отдельного тестового пользователя.
func Register(t *testing.T, request models.RegisterRequest) models.LoginOkResponse {
	response := auth.HttpPostAuthRegister(t, request)
	defer response.Body.Close()

	require.Equal(t, http.StatusCreated, response.StatusCode, "Значение response.StatusCode не совпадает с ожидаемым")

	body, err := io.ReadAll(response.Body)
	require.NoError(t, err, "Операция должна завершиться без ошибки")

	var result models.LoginOkResponse
	require.NoError(t, json.Unmarshal(body, &result), "Ответ API должен соответствовать JSON-модели")

	return result
}
