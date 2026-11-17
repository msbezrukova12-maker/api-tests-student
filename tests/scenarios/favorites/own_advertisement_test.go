package favorites

import (
	"testing"

	"api-tests-template/internal/managers/advertisements"
	manager "api-tests-template/internal/managers/favorites"
	"api-tests-template/internal/managers/favorites/models"
	"api-tests-template/tests/fixtures"
)

// TestFavoritesOwnAdvertisement проверяет попытку добавить собственное объявление — пример из задания Avito.
func TestFavoritesOwnAdvertisement(t *testing.T) {
	_, viewer := fixtures.FavoritesUsers(t)
	advertisement := advertisements.CreateAdvertisement(t, viewer.Token, fixtures.AdvertisementRequest(t))

	t.Log("[Step] Пытаемся добавить собственное объявление в избранное")
	status, response := manager.AddFavoriteExpectClientError(
		t,
		viewer.Token,
		models.AddFavoriteRequest{AdvertisementID: advertisement.ID},
	)
	t.Logf("[Result] Получена клиентская ошибка HTTP %d: %s", status, response.Message)

	if items := manager.GetAllFavorites(t, viewer.Token); len(items) != 0 {
		t.Fatalf("Собственное объявление не должно появляться в избранном: %d записей", len(items))
	}
}
