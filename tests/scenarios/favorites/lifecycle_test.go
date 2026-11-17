package favorites

import (
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/managers/advertisements"
	manager "api-tests-template/internal/managers/favorites"
	"api-tests-template/tests/fixtures"
)

// TestFavoritesLifecycle проверяет добавление, удаление и каскадное удаление после DELETE объявления.
func TestFavoritesLifecycle(t *testing.T) {
	owner, viewer := fixtures.FavoritesUsers(t)

	t.Log("[Precondition] Создаём чужое для viewer объявление")
	request := fixtures.AdvertisementRequest(t)
	advertisement := advertisements.CreateAdvertisement(t, owner.Token, request)

	t.Log("[Step] Добавляем объявление в избранное")
	manager.AddFavorite(t, viewer.Token, advertisement.ID)
	items := manager.GetAllFavorites(t, viewer.Token)
	listed := assertions.FindFavoriteByAdvertisementID(t, items, advertisement.ID)
	require.Equal(t, advertisement.ID, listed.Advertisement.ID, "GET /favorites вернул другое объявление")

	t.Log("[Step] Удаляем объявление из избранного")
	manager.RemoveFavorite(t, viewer.Token, advertisement.ID)
	require.Empty(t, manager.GetAllFavorites(t, viewer.Token),
		"После DELETE /favorites/{id} объявление осталось в избранном")

	t.Log("[Step] Снова добавляем объявление, затем удаляем само объявление")
	manager.AddFavorite(t, viewer.Token, advertisement.ID)
	advertisements.DeleteAdvertisement(t, owner.Token, advertisement.ID)
	require.Empty(t, manager.GetAllFavorites(t, viewer.Token),
		"После DELETE объявления оно осталось в избранном")

	t.Log("[Result] Полный цикл favorites выполнен успешно")
}
