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

func HttpGetAdvertisements(t *testing.T, token, textForSearch string) *http.Response {
	t.Helper()

	return apiRunner.GetRunner().Auth(token).Create().
		Get(path.AdvertisementsPath).
		Query("search", textForSearch).
		Expect(t).
		Status(http.StatusOK).
		End().Response
}

func HttpPostAdvertisement(t *testing.T, token string, request models.CreateAdvertisementRequest) *http.Response {
	t.Helper()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	for name, value := range request.Fields {
		require.NoError(t, writer.WriteField(name, value))
	}

	for _, photo := range request.Photos {
		part, err := writer.CreateFormFile("photos", photo.Name)
		require.NoError(t, err)

		_, err = part.Write(photo.Content)
		require.NoError(t, err)
	}

	require.NoError(t, writer.Close())

	return apiRunner.GetRunner().Auth(token).Create().
		Post(path.AdvertisementPath).
		ContentType(writer.FormDataContentType()).
		Body(body.String()).
		Expect(t).
		End().Response
}

func HttpGetAdvertisement(t *testing.T, token, id string) *http.Response {
	t.Helper()

	return apiRunner.GetRunner().Auth(token).Create().
		Get(path.AdvertisementPath).
		Query("id", id).
		Expect(t).
		End().Response
}

func HttpDeleteAdvertisement(t *testing.T, token, id string) *http.Response {
	t.Helper()

	return apiRunner.GetRunner().Auth(token).Create().
		Delete(path.AdvertisementPath).
		Query("id", id).
		Expect(t).
		End().Response
}

func HttpGetAdvertisementPhotos(t *testing.T, token, id string) *http.Response {
	t.Helper()

	return apiRunner.GetRunner().Auth(token).Create().
		Get(fmt.Sprintf("%s/%s/photos", path.AdvertisementsPath, id)).
		Expect(t).
		End().Response
}
