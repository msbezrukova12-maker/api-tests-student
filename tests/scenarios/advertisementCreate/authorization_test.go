package advertisementCreate

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/tests/fixtures"
)

const invalidToken = "incorrect_token"

func TestCreateAdvertisementWithoutAuthorization(t *testing.T) {
	cases := []struct {
		name  string
		token string
	}{
		{name: "missing_token"},
		{name: "invalid_token", token: invalidToken},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fixtures.LoadEnvironment(t)
			request := fixtures.AdvertisementRequest(t)

			response := advertisements.CreateAdvertisementExpectError(
				t,
				tc.token,
				request,
				http.StatusUnauthorized,
			)

			require.NotEmpty(t, response.Error)
			require.NotEmpty(t, response.Message)
		})
	}
}
