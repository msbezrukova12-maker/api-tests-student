package favorites

import (
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/advertisements"
	manager "api-tests-template/internal/managers/favorites"
	"api-tests-template/internal/managers/favorites/models"
	"api-tests-template/tests/fixtures"
)

// TestFavoritesDuplicateAdvertisement проверяет повторное добавление одного объявления — пример из задания Avito.
func TestFavoritesDuplicateAdvertisement(t *testing.T) {
	owner, viewer := fixtures.FavoritesUsers(t)
	advertisement := advertisements.CreateAdvertisement(t, owner.Token, fixtures.AdvertisementRequest(t))

	t.Log("[Precondition] Добавляем чужое объявление в избранное")
	saved := manager.AddFavorite(t, viewer.Token, advertisement.ID)

	t.Log("[Step] Повторно добавляем то же объявление")
	status, response := manager.AddFavoriteExpectClientError(
		t,
		viewer.Token,
		models.AddFavoriteRequest{AdvertisementID: advertisement.ID},
	)
	t.Logf("[Check] Повторный POST вернул HTTP %d: %s", status, response.Message)

	items := manager.GetAllFavorites(t, viewer.Token)
	require.Len(t, items, 1, "Повторный POST не должен создавать дубликат")
	require.Equal(t, saved.FavoriteID, items[0].FavoriteID, "Исходная запись favorite изменилась")
}
