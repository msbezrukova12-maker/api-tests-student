package assertions

import (
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/favorites/models"
)

// FindFavoriteByAdvertisementID возвращает запись избранного для указанного объявления.
func FindFavoriteByAdvertisementID(
	t *testing.T,
	items []models.Favorite,
	advertisementID string,
) models.Favorite {
	t.Helper()

	for _, item := range items {
		if item.Advertisement.ID == advertisementID {
			return item
		}
	}

	require.FailNow(t, "Favorite не найден",
		"Ожидали advertisement_id=%s среди %d записей", advertisementID, len(items))
	return models.Favorite{}
}
