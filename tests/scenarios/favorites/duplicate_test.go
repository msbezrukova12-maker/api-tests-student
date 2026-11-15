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

// TestFavoritesDuplicateAdvertisement проверяет запрет повторного добавления одного объявления.
func TestFavoritesDuplicateAdvertisement(t *testing.T) {
	owner, viewer := fixtures.FavoritesUsers(t)
	advertisement := advertisements.CreateAdvertisement(t, owner.Token, fixtures.AdvertisementRequest(t))

	t.Log("[Precondition] Один раз добавляем чужое объявление в избранное")
	saved := manager.AddFavorite(t, viewer.Token, advertisement.ID)

	t.Log("[Step] Повторно добавляем то же объявление")
	response := manager.AddFavoriteExpectError(
		t,
		viewer.Token,
		models.AddFavoriteRequest{AdvertisementID: advertisement.ID},
		http.StatusConflict,
	)
	require.Equal(t, "already_exists", response.Error, "Неправильный duplicate error")
	require.Equal(t, "Advertisement is already in favorites", response.Message, "Неправильное duplicate message")

	items := manager.GetAllFavorites(t, viewer.Token)
	require.Len(t, items, 1, "Повторный POST создал дубликат")
	require.Equal(t, saved.FavoriteID, items[0].FavoriteID, "Повторный POST изменил исходный favorite")
	t.Log("[Result] Повторное добавление отклонено HTTP 409")
}
