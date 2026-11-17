package favorites

import (
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/advertisements"
	manager "api-tests-template/internal/managers/favorites"
	"api-tests-template/internal/managers/favorites/models"
	"api-tests-template/tests/fixtures"
)

func TestFavoritesNonexistentAdvertisement(t *testing.T) {
	owner, viewer := fixtures.FavoritesUsers(t)

	advertisement := advertisements.CreateAdvertisement(t, owner.Token, fixtures.AdvertisementRequest(t))
	advertisements.DeleteAdvertisement(t, owner.Token, advertisement.ID)

	_, addError := manager.AddFavoriteExpectClientError(
		t,
		viewer.Token,
		models.AddFavoriteRequest{AdvertisementID: advertisement.ID},
	)
	require.NotEmpty(t, addError.Message)

	_, removeError := manager.RemoveFavoriteExpectClientError(t, viewer.Token, advertisement.ID)
	require.NotEmpty(t, removeError.Message)

	require.Empty(t, manager.GetAllFavorites(t, viewer.Token))
}
