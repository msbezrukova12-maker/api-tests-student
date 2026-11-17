package favorites

import (
	"net/http"
	"testing"

	manager "api-tests-template/internal/managers/favorites"
	"api-tests-template/tests/fixtures"
)

func TestFavoritesMissingAdvertisementID(t *testing.T) {
	t.Log("Precondition: авторизуем пользователя")
	user := fixtures.Login(t)
	t.Log("Step: добавляем в избранное без обязательного advertisement_id")
	response := manager.AddFavoriteExpectStatus(t, user.Token, map[string]string{}, http.StatusBadRequest)
	fixtures.AssertError(t, response, fixtures.InvalidRequestError,
		"Key: 'AddFavoriteRequest.AdvertisementID' Error:Field validation for 'AdvertisementID' failed on the 'required' tag")
	t.Log("Result: получена точная ошибка обязательного поля")
}
