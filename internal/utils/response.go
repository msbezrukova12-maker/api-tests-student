package utils

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func ReadResponseBody(t *testing.T, response *http.Response) []byte {
	t.Helper()

	require.NotNil(t, response)
	require.NotNil(t, response.Body)
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	require.NoError(t, err)

	return body
}

func DecodeResponse[T any](t *testing.T, response *http.Response, expectedStatus int) T {
	t.Helper()

	body := ReadResponseBody(t, response)
	require.Equal(t, expectedStatus, response.StatusCode, "response: %s", body)

	var result T
	require.NoError(t, json.Unmarshal(body, &result))

	return result
}
