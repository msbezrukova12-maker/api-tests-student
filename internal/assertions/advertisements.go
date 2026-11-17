package assertions

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/advertisements/models"
)

// ContainsAdvertisementByTitle проверяет наличие полного названия в поиске.
func ContainsAdvertisementByTitle(t *testing.T, items []models.Advertisement, title string) {
	t.Helper()

	for _, item := range items {
		if item.Title == title {
			return
		}
	}

	require.Failf(t, "Объявление не найдено", "title=%q", title)
}

// AdvertisementMatchesRequest сравнивает поля объявления с отправленным запросом.
func AdvertisementMatchesRequest(
	t *testing.T,
	actual models.Advertisement,
	request models.CreateAdvertisementRequest,
	userID string,
) {
	t.Helper()

	require.NotEmpty(t, actual.ID, "Поле actual.ID не должно быть пустым")
	require.Equal(t, request.Fields["title"], actual.Title, "Сохранённый title отличается от отправленного")
	require.Equal(t, request.Fields["description"], actual.Description, "Сохранённый description отличается от отправленного")

	if value, ok := request.Fields["price"]; ok {
		price, err := strconv.ParseFloat(value, 64)
		require.NoError(t, err, "Операция должна завершиться без ошибки")
		require.Equal(t, price, actual.Price, "Значение actual.Price не совпадает с ожидаемым")
	}

	if value, ok := request.Fields["quantity"]; ok {
		quantity, err := strconv.Atoi(value)
		require.NoError(t, err, "Операция должна завершиться без ошибки")
		require.Equal(t, quantity, actual.Quantity, "Значение actual.Quantity не совпадает с ожидаемым")
	}

	require.Equal(t, userID, actual.UserID, "Значение actual.UserID не совпадает с ожидаемым")
	require.NotEmpty(t, actual.CreatedAt, "Поле actual.CreatedAt не должно быть пустым")
	require.NotEmpty(t, actual.UpdatedAt, "Поле actual.UpdatedAt не должно быть пустым")
	require.Len(t, actual.Photos, len(request.Photos), "Число фотографий объявления отличается от отправленного")
}
