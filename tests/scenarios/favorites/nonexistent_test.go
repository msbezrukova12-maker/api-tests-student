package favorites

import (
	"testing"

	"api-tests-template/internal/managers/advertisements"
	manager "api-tests-template/internal/managers/favorites"
	"api-tests-template/internal/managers/favorites/models"
	"api-tests-template/tests/fixtures"
)

// TestFavoritesNonexistentAdvertisement проверяет добавление и удаление несуществующего ad_id из задания Avito.
func TestFavoritesNonexistentAdvertisement(t *testing.T) {
	owner, viewer := fixtures.FavoritesUsers(t)

	t.Log("[Precondition] Создаём и удаляем объявление, чтобы получить гарантированно несуществующий ID")
	advertisement := advertisements.CreateAdvertisement(t, owner.Token, fixtures.AdvertisementRequest(t))
	advertisements.DeleteAdvertisement(t, owner.Token, advertisement.ID)

	t.Run("add", func(t *testing.T) {
		t.Log("[Step] Добавляем несуществующий ad_id в избранное")
		status, response := manager.AddFavoriteExpectClientError(
			t,
			viewer.Token,
			models.AddFavoriteRequest{AdvertisementID: advertisement.ID},
		)
		t.Logf("[Result] Получена клиентская ошибка HTTP %d: %s", status, response.Message)
	})

	t.Run("remove", func(t *testing.T) {
		t.Log("[Step] Удаляем несуществующий ad_id из избранного")
		status, response := manager.RemoveFavoriteExpectClientError(t, viewer.Token, advertisement.ID)
		t.Logf("[Result] Получена клиентская ошибка HTTP %d: %s", status, response.Message)
	})

	if items := manager.GetAllFavorites(t, viewer.Token); len(items) != 0 {
		t.Fatalf("Негативные запросы изменили favorites: %d записей", len(items))
	}
}
