package favorites

import (
	"net/http"
	"testing"
	"time"

	"api-tests-template/internal/managers/advertisements"
	manager "api-tests-template/internal/managers/favorites"
	"api-tests-template/internal/managers/favorites/models"
	"api-tests-template/tests/fixtures"
)

const (
	favoritesLimit       = 50
	limitRequestInterval = time.Second
)

func TestFavoritesLimit(t *testing.T) {
	t.Log("Precondition: готовим владельца объявлений и пользователя с пустым избранным")
	owner, viewer := fixtures.FavoritesUsers(t)

	t.Log("Step: добавляем в избранное 50 отдельных объявлений")
	for i := 0; i < favoritesLimit; i++ {
		time.Sleep(limitRequestInterval)
		advertisement := advertisements.CreateAdvertisement(t, owner.Token, fixtures.AdvertisementRequest(t))
		manager.AddFavorite(t, viewer.Token, advertisement.ID)
		// Пауза перед очисткой, чтобы не отправлять сотню запросов подряд.
		t.Cleanup(func() { time.Sleep(limitRequestInterval) })
	}

	t.Log("Step: пытаемся добавить 51-е объявление")
	advertisement := advertisements.CreateAdvertisement(t, owner.Token, fixtures.AdvertisementRequest(t))
	response := manager.AddFavoriteExpectStatus(t, viewer.Token,
		models.AddFavoriteRequest{AdvertisementID: advertisement.ID}, http.StatusBadRequest)
	fixtures.AssertError(t, response, "limit_reached", "Favorites limit of 50 reached")
	t.Log("Result: 51-е объявление не добавлено, получена точная ошибка лимита")
}
