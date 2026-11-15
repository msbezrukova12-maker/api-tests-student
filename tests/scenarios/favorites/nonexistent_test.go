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

// TestFavoritesNonexistentAdvertisement проверяет добавление и удаление несуществующего ad_id.
func TestFavoritesNonexistentAdvertisement(t *testing.T) {
	owner, viewer := fixtures.FavoritesUsers(t)

	t.Log("[Precondition] Создаём и удаляем объявление, чтобы получить гарантированно несуществующий UUID")
	advertisement := advertisements.CreateAdvertisement(t, owner.Token, fixtures.AdvertisementRequest(t))
	advertisements.DeleteAdvertisement(t, owner.Token, advertisement.ID)

	t.Run("add", func(t *testing.T) {
		t.Log("[Step] Добавляем удалённое объявление в избранное")
		response := manager.AddFavoriteExpectError(
			t,
			viewer.Token,
			models.AddFavoriteRequest{AdvertisementID: advertisement.ID},
			http.StatusNotFound,
		)
		require.Equal(t, "not_found", response.Error, "Неправильный error для отсутствующего объявления")
		require.Equal(t, "Advertisement not found", response.Message, "Неправильное message")
	})

	t.Run("remove", func(t *testing.T) {
		t.Log("[Step] Удаляем отсутствующую запись избранного")
		response := manager.RemoveFavoriteExpectError(
			t,
			viewer.Token,
			advertisement.ID,
			http.StatusNotFound,
		)
		require.Equal(t, "not_found", response.Error, "Неправильный error для отсутствующего favorite")
		require.Equal(t, "Favorite entry not found", response.Message, "Неправильное message")
	})

	require.Empty(t, manager.GetAllFavorites(t, viewer.Token), "Негативные запросы изменили favorites")
	t.Log("[Result] Обе операции с несуществующим ad_id вернули HTTP 404")
}
