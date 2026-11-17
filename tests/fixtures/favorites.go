package fixtures

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/auth"
	authModels "api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/utils"
)

// FavoritesUsers готовит владельца объявления и отдельного пользователя избранного.
func FavoritesUsers(t *testing.T) (authModels.LoginOkResponse, authModels.LoginOkResponse) {
	t.Helper()

	owner := Login(t)

	viewer := auth.Register(t, authModels.RegisterRequest{
		Email:     strings.ToLower(utils.RandomString(12)) + "@example.com",
		Password:  "Aa1!" + utils.RandomString(12),
		FirstName: "Test",
	})

	require.NotEmpty(t, viewer.Token, "Поле viewer.Token не должно быть пустым")
	require.NotEmpty(t, viewer.User.Id, "Поле viewer.User.Id не должно быть пустым")
	require.NotEqual(t, owner.User.Id, viewer.User.Id, "Владелец и пользователь избранного должны различаться")

	return owner, viewer
}
