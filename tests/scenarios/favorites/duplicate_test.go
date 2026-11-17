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

func TestFavoritesDuplicateAdvertisement(t *testing.T) {
	t.Log("Precondition: подготовка пользователя и тестовых данных")
	owner, viewer := fixtures.FavoritesUsers(t)
	advertisement := advertisements.CreateAdvertisement(t, owner.Token, fixtures.AdvertisementRequest(t))

	t.Log("Step: добавляем объявление и повторяем добавление")
	saved := manager.AddFavorite(t, viewer.Token, advertisement.ID)

	response := manager.AddFavoriteExpectStatus(
		t,
		viewer.Token,
		models.AddFavoriteRequest{AdvertisementID: advertisement.ID},
		http.StatusConflict,
	)
	fixtures.AssertError(t, response, "already_exists", "Advertisement is already in favorites")

	items := manager.GetAllFavorites(t, viewer.Token)
	require.Len(t, items, 1, "Повторное добавление не должно создавать новую запись")
	require.Equal(t, saved.FavoriteID, items[0].FavoriteID, "Повторное добавление изменило запись избранного")
	t.Log("Result: проверки сценария пройдены")
}
