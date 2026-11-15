package myAdvertisement

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"

	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/managers/myAdvertisements"
	"api-tests-template/tests/fixtures"
)

// TestGetMyAdvertisementsPositive проверяет, что созданное объявление появляется в списке текущего пользователя.
func TestGetMyAdvertisementsPositive(t *testing.T) {
	user := fixtures.Login(t)

	t.Log("[Precondition] Создаём независимое объявление текущего пользователя")
	created := advertisements.CreateAdvertisement(t, user.Token, fixtures.AdvertisementRequest(t))

	t.Log("[Step] Получаем список собственных объявлений")
	body := myAdvertisements.GetMyAdvertisements(t, user.Token, http.StatusOK, 50)
	items := gjson.Get(body, "items").Array()
	require.NotEmpty(t, items, "Список собственных объявлений пуст после успешного создания")

	found := false
	for _, item := range items {
		if gjson.Get(item.String(), "id").String() != created.ID {
			continue
		}

		found = true
		require.Equal(t, user.User.Id, gjson.Get(item.String(), "user_id").String(), "user_id не совпадает")
		require.NotEmpty(t, gjson.Get(item.String(), "created_at").String(), "created_at пустое")
		require.NotEmpty(t, gjson.Get(item.String(), "updated_at").String(), "updated_at пустое")
		require.NotEmpty(t, gjson.Get(item.String(), "title").String(), "Название пустое")
		require.NotEmpty(t, gjson.Get(item.String(), "description").String(), "Описание пустое")
		require.NotEmpty(t, gjson.Get(item.String(), "photos").Array(), "Фотографии пустые")
		break
	}

	require.True(t, found, "Созданное объявление %s отсутствует в /my/advertisements", created.ID)
	t.Log("[Result] Созданное объявление найдено среди собственных")
}

// TestGetMyAdvertisementsIncorrectToken проверяет HTTP 401 и сообщение для невалидного токена.
func TestGetMyAdvertisementsIncorrectToken(t *testing.T) {
	fixtures.LoadEnvironment(t)

	t.Log("[Step] Получаем собственные объявления с неправильным токеном")
	body := myAdvertisements.GetMyAdvertisements(t, "incorrect_token", http.StatusUnauthorized, 50)

	require.Equal(t, "unauthorized", gjson.Get(body, "error").String(), "Некорректный error")
	require.Equal(t, "Invalid or expired token", gjson.Get(body, "message").String(), "Некорректный message")
	t.Log("[Result] Получен ожидаемый HTTP 401")
}
