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

func TestFavoritesNonexistentAdvertisement(t *testing.T) {
	t.Log("Precondition: подготовка пользователя и тестовых данных")
	owner, viewer := fixtures.FavoritesUsers(t)

	advertisement := advertisements.CreateAdvertisement(t, owner.Token, fixtures.AdvertisementRequest(t))
	advertisements.DeleteAdvertisement(t, owner.Token, advertisement.ID)

	t.Log("Step: добавляем и удаляем несуществующий ID")
	addError := manager.AddFavoriteExpectStatus(
		t,
		viewer.Token,
		models.AddFavoriteRequest{AdvertisementID: advertisement.ID},
		http.StatusNotFound,
	)
	fixtures.AssertError(t, addError, "not_found", "Advertisement not found")

	removeError := manager.RemoveFavoriteExpectStatus(t, viewer.Token, advertisement.ID, http.StatusNotFound)
	fixtures.AssertError(t, removeError, "not_found", "Favorite entry not found")

	require.Empty(t, manager.GetAllFavorites(t, viewer.Token), "Избранное должно быть пустым после операции")
	t.Log("Result: проверки сценария пройдены")
}
