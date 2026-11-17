package fixtures

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/auth"
	authModels "api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/utils"
)

func FavoritesUsers(t *testing.T) (authModels.LoginOkResponse, authModels.LoginOkResponse) {
	t.Helper()

	owner := Login(t)

	viewer := auth.Register(t, authModels.RegisterRequest{
		Email:     strings.ToLower(utils.RandomString(12)) + "@example.com",
		Password:  "Aa1!" + utils.RandomString(12),
		FirstName: "Test",
	})

	require.NotEmpty(t, viewer.Token)
	require.NotEmpty(t, viewer.User.Id)
	require.NotEqual(t, owner.User.Id, viewer.User.Id)

	return owner, viewer
}
