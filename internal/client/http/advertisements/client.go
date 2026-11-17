package advertisements

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
	"api-tests-template/internal/managers/advertisements/models"
)

// HttpGetAdvertisements выполняет GET-запрос списка объявлений с фильтрацией по названию.
func HttpGetAdvertisements(t *testing.T, token, textForSearch string) *http.Response {
	t.Helper()
	return apiRunner.GetRunner().Auth(token).Create().Get(path.AdvertisementsPath).
		Query("search", textForSearch).
		Expect(t).
		Status(http.StatusOK).
		End().Response
}

// HttpPostAdvertisement отправляет multipart-запрос создания объявления.
func HttpPostAdvertisement(t *testing.T, token string, request models.CreateAdvertisementRequest) *http.Response {
	t.Helper()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	for name, value := range request.Fields {
		require.NoError(t, writer.WriteField(name, value), "Не удалось записать multipart field %q", name)
	}

	for _, photo := range request.Photos {
		part, err := writer.CreateFormFile("photos", photo.Name)
		require.NoError(t, err, "Не удалось создать multipart photo %q", photo.Name)

		_, err = part.Write(photo.Content)
		require.NoError(t, err, "Не удалось записать photo %q", photo.Name)
	}

	require.NoError(t, writer.Close(), "Не удалось завершить multipart request")

	return apiRunner.GetRunner().Auth(token).Create().Post(path.AdvertisementPath).
		ContentType(writer.FormDataContentType()).
		Body(body.String()).
		Expect(t).
		End().Response
}

// HttpGetAdvertisement получает одно объявление по ID.
func HttpGetAdvertisement(t *testing.T, token, id string) *http.Response {
	t.Helper()
	return apiRunner.GetRunner().Auth(token).Create().Get(path.AdvertisementPath).
		Query("id", id).
		Expect(t).
		End().Response
}

// HttpDeleteAdvertisement удаляет объявление по ID.
func HttpDeleteAdvertisement(t *testing.T, token, id string) *http.Response {
	t.Helper()
	return apiRunner.GetRunner().Auth(token).Create().Delete(path.AdvertisementPath).
		Query("id", id).
		Expect(t).
		End().Response
}

// HttpGetAdvertisementPhotos получает список фотографий объявления.
func HttpGetAdvertisementPhotos(t *testing.T, token, id string) *http.Response {
	t.Helper()
	return apiRunner.GetRunner().Auth(token).Create().
		Get(fmt.Sprintf("%s/%s/photos", path.AdvertisementsPath, id)).
		Expect(t).
		End().Response
}
