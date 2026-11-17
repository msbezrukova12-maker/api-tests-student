package fixtures

import (
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/advertisements/models"
)

const (
	InvalidToken                = "incorrect_token"
	UnauthorizedError           = "unauthorized"
	MissingTokenMessage         = "Authorization header required"
	InvalidAuthorizationMessage = "Invalid or expired token"
	InvalidRequestError         = "invalid_request"
	RequiredTitleMessage        = "title is required"
	PriceValidationMessage      = "price must be between 0 and 1000000 inclusive"
)

// AssertError проверяет точные error и message ответа API.
func AssertError(t *testing.T, actual models.ErrorResponse, code, message string) {
	t.Helper()
	require.Equal(t, code, actual.Error, "Некорректный error в ответе API")
	require.Equal(t, message, actual.Message, "Некорректный message в ответе API")
}
