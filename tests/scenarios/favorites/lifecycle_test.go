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
	t.Log("Precondition: подготовка пользователя и тестовых данных")
	owner, viewer := fixtures.FavoritesUsers(t)
	advertisement := advertisements.CreateAdvertisement(t, owner.Token, fixtures.AdvertisementRequest(t))

	t.Log("Step: добавляем объявление в избранное")
	manager.AddFavorite(t, viewer.Token, advertisement.ID)
	items := manager.GetAllFavorites(t, viewer.Token)
	favorite := assertions.FindFavoriteByAdvertisementID(t, items, advertisement.ID)
	require.Equal(t, advertisement.ID, favorite.Advertisement.ID, "Значение favorite.Advertisement.ID не совпадает с ожидаемым")

	t.Log("Step: удаляем из избранного и проверяем отсутствие записи")
	manager.RemoveFavorite(t, viewer.Token, advertisement.ID)
	require.Empty(t, manager.GetAllFavorites(t, viewer.Token), "Избранное должно быть пустым после операции")

	t.Log("Step: добавляем объявление в избранное")
	manager.AddFavorite(t, viewer.Token, advertisement.ID)
	t.Log("Step: удаляем само объявление и проверяем избранное")
	advertisements.DeleteAdvertisement(t, owner.Token, advertisement.ID)
	require.Empty(t, manager.GetAllFavorites(t, viewer.Token), "Избранное должно быть пустым после операции")
	t.Log("Result: проверки сценария пройдены")
}
