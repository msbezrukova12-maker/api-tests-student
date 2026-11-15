package advertisementCreate

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/tests/fixtures"
)

// TestCreateAdvertisementAllFields проверяет создание объявления со всеми параметрами, поиск, GET по ID и фотографии.
func TestCreateAdvertisementAllFields(t *testing.T) {
	user := fixtures.Login(t)

	t.Log("[Precondition] Подготавливаем объявление со всеми параметрами и тремя фотографиями")
	request := fixtures.AdvertisementRequest(t)
	request.Photos = append(request.Photos,
		fixtures.PNGPhoto(t, "photo-2.png", color.RGBA{R: 210, G: 80, B: 90, A: 255}),
		fixtures.PNGPhoto(t, "photo-3.png", color.RGBA{R: 60, G: 100, B: 220, A: 255}),
	)

	t.Log("[Step] Создаём объявление")
	created := advertisements.CreateAdvertisement(t, user.Token, request)
	assertions.AdvertisementMatchesRequest(t, created, request, user.User.Id)

	t.Log("[Step] Получаем созданное объявление по ID")
	actual := advertisements.GetAdvertisement(t, user.Token, created.ID)
	require.Equal(t, created.ID, actual.ID, "GET вернул другое объявление")
	assertions.AdvertisementMatchesRequest(t, actual, request, user.User.Id)

	t.Log("[Step] Ищем объявление по полному title")
	list := advertisements.GetAdvertisementsList(t, user.Token, request.Fields["title"])
	found := assertions.FindAdvertisementByID(t, list.Items, created.ID)
	assertions.AdvertisementMatchesRequest(t, found, request, user.User.Id)

	t.Log("[Step] Получаем фотографии через GET /advertisements/{id}/photos")
	photos := advertisements.GetAdvertisementPhotos(t, user.Token, created.ID)
	require.Len(t, photos, len(request.Photos), "Photo endpoint вернул неправильное количество фотографий")

	for i, photo := range photos {
		require.NotEmpty(t, photo.ID, "Photo %d не содержит ID", i)
		require.Equal(t, created.ID, photo.AdvertisementID, "Photo %d относится к другому объявлению", i)
		require.NotEmpty(t, photo.URL, "Photo %d не содержит URL", i)
		require.Equal(t, i, photo.SortOrder, "Photo %d имеет неправильный sort_order", i)
		require.NotEmpty(t, photo.CreatedAt, "Photo %d не содержит created_at", i)
	}

	t.Log("[Result] Объявление создано, читается по ID, находится по полному title, фотографии доступны")
}
