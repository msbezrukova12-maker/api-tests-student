package favorites

import (
	"net/http"
	"testing"

	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
)

// HttpGetFavorites получает избранное пользователя.
func HttpGetFavorites(t *testing.T, token string) *http.Response {
	t.Helper()

	return apiRunner.GetRunner().Auth(token).Create().
		Get(path.FavoritesPath).
		Expect(t).
		End().Response
}

// HttpPostFavorite отправляет запрос добавления в избранное.
func HttpPostFavorite(t *testing.T, token string, request any) *http.Response {
	t.Helper()

	return apiRunner.GetRunner().Auth(token).Create().
		Post(path.FavoritesPath).
		ContentType("application/json").
		JSON(request).
		Expect(t).
		End().Response
}

// HttpDeleteFavorite удаляет объявление из избранного.
func HttpDeleteFavorite(t *testing.T, token, advertisementID string) *http.Response {
	t.Helper()

	return apiRunner.GetRunner().Auth(token).Create().
		Delete(path.FavoritesPath + "/" + advertisementID).
		Expect(t).
		End().Response
}
