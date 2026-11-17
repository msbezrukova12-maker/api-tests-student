package assertions

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/advertisements/models"
)

func ContainsAdvertisementByTitle(t *testing.T, items []models.Advertisement, title string) {
	t.Helper()

	for _, item := range items {
		if item.Title == title {
			return
		}
	}

	require.Failf(t, "Объявление не найдено", "title=%q", title)
}

func FindAdvertisementByID(t *testing.T, items []models.Advertisement, id string) models.Advertisement {
	t.Helper()

	for _, item := range items {
		if item.ID == id {
			return item
		}
	}

	require.FailNow(t, "Объявление не найдено", "id=%s", id)
	return models.Advertisement{}
}

func AdvertisementMatchesRequest(
	t *testing.T,
	actual models.Advertisement,
	request models.CreateAdvertisementRequest,
	userID string,
) {
	t.Helper()

	require.NotEmpty(t, actual.ID)
	require.Equal(t, request.Fields["title"], actual.Title)
	require.Equal(t, request.Fields["description"], actual.Description)

	if value, ok := request.Fields["price"]; ok {
		price, err := strconv.ParseFloat(value, 64)
		require.NoError(t, err)
		require.Equal(t, price, actual.Price)
	}

	if value, ok := request.Fields["quantity"]; ok {
		quantity, err := strconv.Atoi(value)
		require.NoError(t, err)
		require.Equal(t, quantity, actual.Quantity)
	}

	require.Equal(t, userID, actual.UserID)
	require.NotEmpty(t, actual.CreatedAt)
	require.NotEmpty(t, actual.UpdatedAt)
	require.Len(t, actual.Photos, len(request.Photos))
}
