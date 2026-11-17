package favorites

import (
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/managers/advertisements"
	manager "api-tests-template/internal/managers/favorites"
	"api-tests-template/tests/fixtures"
)

func TestFavoritesLifecycle(t *testing.T) {
	owner, viewer := fixtures.FavoritesUsers(t)
	advertisement := advertisements.CreateAdvertisement(t, owner.Token, fixtures.AdvertisementRequest(t))

	manager.AddFavorite(t, viewer.Token, advertisement.ID)
	items := manager.GetAllFavorites(t, viewer.Token)
	favorite := assertions.FindFavoriteByAdvertisementID(t, items, advertisement.ID)
	require.Equal(t, advertisement.ID, favorite.Advertisement.ID)

	manager.RemoveFavorite(t, viewer.Token, advertisement.ID)
	require.Empty(t, manager.GetAllFavorites(t, viewer.Token))

	manager.AddFavorite(t, viewer.Token, advertisement.ID)
	advertisements.DeleteAdvertisement(t, owner.Token, advertisement.ID)
	require.Empty(t, manager.GetAllFavorites(t, viewer.Token))
}
