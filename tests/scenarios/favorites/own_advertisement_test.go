package favorites

import (
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/advertisements"
	manager "api-tests-template/internal/managers/favorites"
	"api-tests-template/internal/managers/favorites/models"
	"api-tests-template/tests/fixtures"
)

func TestFavoritesOwnAdvertisement(t *testing.T) {
	_, viewer := fixtures.FavoritesUsers(t)
	advertisement := advertisements.CreateAdvertisement(t, viewer.Token, fixtures.AdvertisementRequest(t))

	_, response := manager.AddFavoriteExpectClientError(
		t,
		viewer.Token,
		models.AddFavoriteRequest{AdvertisementID: advertisement.ID},
	)

	require.NotEmpty(t, response.Message)
	require.Empty(t, manager.GetAllFavorites(t, viewer.Token))
}
