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

const pageLimit = 50

func TestGetMyAdvertisements(t *testing.T) {
	t.Log("Precondition: создаём собственное объявление")
	user := fixtures.Login(t)
	created := advertisements.CreateAdvertisement(t, user.Token, fixtures.AdvertisementRequest(t))

	t.Log("Step: получаем список собственных объявлений")
	body := myAdvertisements.GetMyAdvertisements(t, user.Token, http.StatusOK, pageLimit)
	ids := gjson.Get(body, "items.#.id").Array()
	found := false
	for _, id := range ids {
		if id.String() == created.ID {
			found = true
			break
		}
	}
	require.True(t, found, "Созданное объявление %s отсутствует в собственных объявлениях", created.ID)
	t.Log("Result: созданное объявление найдено в профиле")
	t.Log("Result: проверки сценария пройдены")
}

func TestGetMyAdvertisementsIncorrectToken(t *testing.T) {
	t.Log("Precondition: загружаем окружение и используем неверный токен")
	fixtures.LoadEnvironment(t)
	t.Log("Step: получаем собственные объявления без действующей авторизации")
	body := myAdvertisements.GetMyAdvertisements(t, fixtures.InvalidToken, http.StatusUnauthorized, pageLimit)
	require.Equal(t, fixtures.UnauthorizedError, gjson.Get(body, "error").String(), "Некорректный error авторизации")
	require.Equal(t, fixtures.InvalidAuthorizationMessage, gjson.Get(body, "message").String(), "Некорректный message авторизации")
	t.Log("Result: получена точная ошибка авторизации")
	t.Log("Result: проверки сценария пройдены")
}
