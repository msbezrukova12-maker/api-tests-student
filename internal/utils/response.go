package utils

import (
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
