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
	user := fixtures.Login(t)

	request := fixtures.AdvertisementRequest(t)
	request.Photos = append(request.Photos,
		fixtures.PNGPhoto(t, "photo-2.png", color.RGBA{R: 210, G: 80, B: 90, A: 255}),
		fixtures.PNGPhoto(t, "photo-3.png", color.RGBA{R: 60, G: 100, B: 220, A: 255}),
	)

	created := advertisements.CreateAdvertisement(t, user.Token, request)
	assertions.AdvertisementMatchesRequest(t, created, request, user.User.Id)

	actual := advertisements.GetAdvertisement(t, user.Token, created.ID)
	require.Equal(t, created.ID, actual.ID)
	assertions.AdvertisementMatchesRequest(t, actual, request, user.User.Id)

	list := advertisements.GetAdvertisementsList(t, user.Token, request.Fields["title"])
	assertions.ContainsAdvertisementByTitle(t, list.Items, request.Fields["title"])

	photos := advertisements.GetAdvertisementPhotos(t, user.Token, created.ID)
	require.Len(t, photos, len(request.Photos))
}
