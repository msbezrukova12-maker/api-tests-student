package favorites

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/advertisements"
	manager "api-tests-template/internal/managers/favorites"
	"api-tests-template/internal/managers/favorites/models"
	"api-tests-template/tests/fixtures"
)

// TestFavoritesOwnAdvertisement проверяет запрет добавления собственного объявления в избранное.
func TestFavoritesOwnAdvertisement(t *testing.T) {
	_, viewer := fixtures.FavoritesUsers(t)
	advertisement := advertisements.CreateAdvertisement(t, viewer.Token, fixtures.AdvertisementRequest(t))

	t.Log("[Step] Пытаемся добавить собственное объявление в избранное")
	response := manager.AddFavoriteExpectError(
		t,
		viewer.Token,
		models.AddFavoriteRequest{AdvertisementID: advertisement.ID},
		http.StatusForbidden,
	)
	require.Equal(t, "forbidden", response.Error, "Неправильный ownership error")
	require.Equal(t, "Cannot add own advertisement to favorites", response.Message, "Неправильное ownership message")
	require.Empty(t, manager.GetAllFavorites(t, viewer.Token),
		"Собственное объявление появилось в избранном")

	t.Log("[Result] Собственное объявление отклонено HTTP 403")
}
