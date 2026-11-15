package favorites

import (
	"net/http"
	"strconv"
	"testing"

	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
)

// HttpGetFavorites получает страницу избранного в стабильной сортировке.
func HttpGetFavorites(t *testing.T, token string, page, pageSize int) *http.Response {
	t.Helper()
	return apiRunner.GetRunner().Auth(token).Create().Get(path.FavoritesPath).
		Query("page", strconv.Itoa(page)).
		Query("page_size", strconv.Itoa(pageSize)).
		Query("sort_by", "title").
		Query("sort_order", "asc").
		Expect(t).
		End().Response
}

// HttpPostFavorite добавляет объявление в избранное.
func HttpPostFavorite(t *testing.T, token string, request any) *http.Response {
	t.Helper()
	return apiRunner.GetRunner().Auth(token).Create().Post(path.FavoritesPath).
		ContentType("application/json").
		JSON(request).
		Expect(t).
		End().Response
}

// HttpDeleteFavorite удаляет объявление из избранного по ID.
func HttpDeleteFavorite(t *testing.T, token, advertisementID string) *http.Response {
	t.Helper()
	return apiRunner.GetRunner().Auth(token).Create().
		Delete(path.FavoritesPath + "/" + advertisementID).
		Expect(t).
		End().Response
}
