package advertisementCreate

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/tests/fixtures"
)

func TestCreateAdvertisementAllFields(t *testing.T) {
	t.Log("Precondition: подготовка пользователя и тестовых данных")
	user := fixtures.Login(t)

	request := fixtures.AdvertisementRequest(t)
	request.Photos = append(request.Photos,
		fixtures.PNGPhoto(t, "photo-2.png", color.RGBA{R: 210, G: 80, B: 90, A: 255}),
		fixtures.PNGPhoto(t, "photo-3.png", color.RGBA{R: 60, G: 100, B: 220, A: 255}),
	)

	t.Log("Step: создаём объявление со всеми полями и тремя фотографиями")
	created := advertisements.CreateAdvertisement(t, user.Token, request)
	assertions.AdvertisementMatchesRequest(t, created, request, user.User.Id)

	t.Log("Step: сверяем сохранённые параметры через GET")
	actual := advertisements.GetAdvertisement(t, user.Token, created.ID)
	require.Equal(t, created.ID, actual.ID, "Значение actual.ID не совпадает с ожидаемым")
	assertions.AdvertisementMatchesRequest(t, actual, request, user.User.Id)

	t.Log("Step: ищем объявление по полному title")
	list := advertisements.GetAdvertisementsList(t, user.Token, request.Fields["title"])
	assertions.ContainsAdvertisementByTitle(t, list.Items, request.Fields["title"])

	t.Log("Step: проверяем фотографии объявления")
	photos := advertisements.GetAdvertisementPhotos(t, user.Token, created.ID)
	require.Len(t, photos, len(request.Photos), "Число загруженных фотографий не совпадает с отправленным")
	require.Equal(t, actual.Photos, photos, "GET photos должен вернуть фотографии созданного объявления")
	for _, photo := range photos {
		require.NotEmpty(t, photo.ID, "У фотографии отсутствует id")
		require.NotEmpty(t, photo.URL, "У фотографии отсутствует URL")
	}
	t.Log("Result: проверки сценария пройдены")
}
