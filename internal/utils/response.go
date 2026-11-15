package utils

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

// ReadResponseBody читает и закрывает HTTP-ответ, сохраняя статус в лог теста.
func ReadResponseBody(t *testing.T, response *http.Response) []byte {
	t.Helper()
	require.NotNil(t, response, "API не вернул HTTP-ответ")
	require.NotNil(t, response.Body, "HTTP-ответ не содержит body")
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	require.NoError(t, err, "Не удалось прочитать HTTP-ответ")
	t.Logf("[HTTP] status=%d", response.StatusCode)
	return body
}

// DecodeResponse проверяет HTTP-статус и декодирует JSON в указанную модель.
func DecodeResponse[T any](t *testing.T, response *http.Response, expectedStatus int) T {
	t.Helper()
	body := ReadResponseBody(t, response)
	require.Equal(t, expectedStatus, response.StatusCode,
		"Неожиданный HTTP status; response: %s", body)

	var result T
	require.NoError(t, json.Unmarshal(body, &result), "Некорректный JSON response: %s", body)
	return result
}
