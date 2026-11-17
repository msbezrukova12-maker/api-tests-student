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

func TestFavoritesOwnAdvertisement(t *testing.T) {
	t.Log("Precondition: подготовка пользователя и тестовых данных")
	_, viewer := fixtures.FavoritesUsers(t)
	advertisement := advertisements.CreateAdvertisement(t, viewer.Token, fixtures.AdvertisementRequest(t))

	t.Log("Step: пытаемся добавить своё объявление в избранное")
	response := manager.AddFavoriteExpectStatus(
		t,
		viewer.Token,
		models.AddFavoriteRequest{AdvertisementID: advertisement.ID},
		http.StatusForbidden,
	)

	fixtures.AssertError(t, response, "forbidden", "Cannot add own advertisement to favorites")
	require.Empty(t, manager.GetAllFavorites(t, viewer.Token), "Избранное должно быть пустым после операции")
	t.Log("Result: проверки сценария пройдены")
}
