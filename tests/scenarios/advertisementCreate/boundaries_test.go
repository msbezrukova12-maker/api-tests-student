package advertisementCreate

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/tests/fixtures"
)

func TestCreateAdvertisementValidBoundaries(t *testing.T) {
	t.Log("Precondition: авторизуем пользователя")
	user := fixtures.Login(t)
	cases := []struct {
		name  string
		field string
		value string
	}{
		{name: "title_50_characters", field: "title", value: strings.Repeat("а", 50)},
		{name: "quantity_100", field: "quantity", value: "100"},
		{name: "price_not_set", field: "price"},
		{name: "one_photo"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Log("Precondition: готовим объявление с одной фотографией")
			request := fixtures.AdvertisementRequest(t)
			if tc.field != "" {
				if tc.value == "" {
					delete(request.Fields, tc.field)
				} else {
					request.Fields[tc.field] = tc.value
				}
			}

			t.Log("Step: публикуем объявление и получаем его через GET")
			created := advertisements.CreateAdvertisement(t, user.Token, request)
			actual := advertisements.GetAdvertisement(t, user.Token, created.ID)
			assertions.AdvertisementMatchesRequest(t, actual, request, user.User.Id)
			if tc.name == "price_not_set" {
				require.Zero(t, actual.Price, "Если цена не указана, API должен сохранить нулевую цену")
			}
			require.Len(t, advertisements.GetAdvertisementPhotos(t, user.Token, created.ID), 1,
				"Объявление должно содержать одну загруженную фотографию")
			t.Log("Result: объявление опубликовано с переданными параметрами")
		})
	}
}
