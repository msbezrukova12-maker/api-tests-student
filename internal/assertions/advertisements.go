package assertions

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/advertisements/models"
)

// ContainsAdvertisementByTitle проверяет наличие объявления с точным названием.
func ContainsAdvertisementByTitle(t *testing.T, items []models.Advertisement, title string) {
	t.Helper()

	for _, item := range items {
		if item.Title == title {
			return
		}
	}

	require.Failf(t, "Объявление не найдено", "В списке нет объявления с названием %q", title)
}

// FindAdvertisementByID возвращает объявление с ожидаемым ID из списка.
func FindAdvertisementByID(t *testing.T, items []models.Advertisement, id string) models.Advertisement {
	t.Helper()

	for _, item := range items {
		if item.ID == id {
			return item
		}
	}

	require.FailNow(t, "Объявление не найдено",
		"Ожидали ID %s среди %d объявлений", id, len(items))
	return models.Advertisement{}
}

// AdvertisementMatchesRequest сравнивает сохранённое объявление со всеми переданными полями.
func AdvertisementMatchesRequest(
	t *testing.T,
	actual models.Advertisement,
	request models.CreateAdvertisementRequest,
	userID string,
) {
	t.Helper()

	require.NotEmpty(t, actual.ID, "Advertisement ID пустой")
	require.Equal(t, request.Fields["title"], actual.Title, "Title отличается от отправленного")
	require.Equal(t, request.Fields["description"], actual.Description, "Description отличается от отправленного")

	if value, present := request.Fields["price"]; present {
		price, err := strconv.ParseFloat(value, 64)
		require.NoError(t, err, "Некорректная цена в fixture")
		require.Equal(t, price, actual.Price, "Price отличается от отправленного")
	}

	if value, present := request.Fields["quantity"]; present {
		quantity, err := strconv.Atoi(value)
		require.NoError(t, err, "Некорректное quantity в fixture")
		require.Equal(t, quantity, actual.Quantity, "Quantity отличается от отправленного")
	}

	require.Equal(t, userID, actual.UserID, "Объявление принадлежит другому пользователю")
	require.NotEmpty(t, actual.CreatedAt, "CreatedAt пустой")
	require.NotEmpty(t, actual.UpdatedAt, "UpdatedAt пустой")
	require.Len(t, actual.Photos, len(request.Photos), "Неожиданное количество фотографий")
}
