package fixtures

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/auth"
	authModels "api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/managers/favorites"
	"api-tests-template/internal/utils"
)

func FavoritesUsers(t *testing.T) (authModels.LoginOkResponse, authModels.LoginOkResponse) {
	t.Helper()

	owner := Login(t)

	request := authModels.RegisterRequest{
		Email:     strings.ToLower(utils.RandomString(12)) + "@example.com",
		Password:  "TestPassword123!",
		FirstName: "Test",
	}

	viewer := auth.RegisterOrLogin(t, request)
	require.NotEmpty(t, viewer.Token)
	require.NotEmpty(t, viewer.User.Id)
	require.NotEqual(t, owner.User.Id, viewer.User.Id)

	cleanFavorites(t, viewer.Token)

	t.Cleanup(func() {
		cleanFavorites(t, viewer.Token)
	})

	return owner, viewer
}

func cleanFavorites(t *testing.T, token string) {
	t.Helper()

	for _, item := range favorites.GetAllFavorites(t, token) {
		favorites.RemoveFavoriteIfExists(t, token, item.Advertisement.ID)
	}
}
