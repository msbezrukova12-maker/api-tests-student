package favorites

import (
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/advertisements"
	manager "api-tests-template/internal/managers/favorites"
	"api-tests-template/internal/managers/favorites/models"
	"api-tests-template/tests/fixtures"
)

func TestFavoritesDuplicateAdvertisement(t *testing.T) {
	owner, viewer := fixtures.FavoritesUsers(t)
	advertisement := advertisements.CreateAdvertisement(t, owner.Token, fixtures.AdvertisementRequest(t))

	saved := manager.AddFavorite(t, viewer.Token, advertisement.ID)

	_, response := manager.AddFavoriteExpectClientError(
		t,
		viewer.Token,
		models.AddFavoriteRequest{AdvertisementID: advertisement.ID},
	)
	require.NotEmpty(t, response.Message)

	items := manager.GetAllFavorites(t, viewer.Token)
	require.Len(t, items, 1)
	require.Equal(t, saved.FavoriteID, items[0].FavoriteID)
}
