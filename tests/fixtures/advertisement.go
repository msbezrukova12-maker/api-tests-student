package fixtures

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/managers/auth"
	authModels "api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/utils"
	base "api-tests-template/tests"
)

const (
	randomTitleLength  = 16
	defaultDescription = "Объявление для автоматической проверки API"
	defaultPrice       = "1234.56"
	defaultQuantity    = "2"
)

// LoadEnvironment загружает общие настройки API без обязательной авторизации.
func LoadEnvironment(t *testing.T) {
	t.Helper()
	base.SetupSuite()
}

// Login загружает окружение и авторизует пользователя из TEST_LOGIN/TEST_PASSWORD.
func Login(t *testing.T) authModels.LoginOkResponse {
	t.Helper()
	LoadEnvironment(t)

	t.Log("[Precondition] Авторизуем пользователя из переменных окружения")
	user := auth.Login(t, os.Getenv("TEST_LOGIN"), os.Getenv("TEST_PASSWORD"))
	require.NotEmpty(t, user.Token, "Login response не содержит token")
	require.NotEmpty(t, user.User.Id, "Login response не содержит user ID")
	return user
}

// AdvertisementRequest создаёт независимый валидный запрос объявления с уникальным title и PNG-фотографией.
func AdvertisementRequest(t *testing.T) models.CreateAdvertisementRequest {
	t.Helper()
	return models.CreateAdvertisementRequest{
		Fields: map[string]string{
			"title":       "API test " + utils.RandomString(randomTitleLength),
			"description": defaultDescription,
			"price":       defaultPrice,
			"quantity":    defaultQuantity,
		},
		Photos: []models.PhotoUpload{
			PNGPhoto(t, "photo.png", color.RGBA{R: 40, G: 150, B: 80, A: 255}),
		},
	}
}

// PNGPhoto генерирует небольшое PNG-изображение без внешнего fixture-файла.
func PNGPhoto(t *testing.T, name string, pixel color.RGBA) models.PhotoUpload {
	t.Helper()

	picture := image.NewRGBA(image.Rect(0, 0, 2, 2))
	picture.SetRGBA(0, 0, pixel)

	var photo bytes.Buffer
	require.NoError(t, png.Encode(&photo, picture), "Не удалось подготовить PNG fixture")
	return models.PhotoUpload{Name: name, Content: photo.Bytes()}
}
