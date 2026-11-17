package assertions

import (
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/favorites/models"
)

// FindFavoriteByAdvertisementID находит избранное по id объявления.
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

	require.FailNow(t, "Favorite не найден", "advertisement_id=%s", advertisementID)
	return models.Favorite{}
}
