package models

// CreateAdvertisementRequest хранит multipart-поля и фотографии, включая заведомо невалидные значения для негативных тестов.
type CreateAdvertisementRequest struct {
	Fields map[string]string
	Photos []PhotoUpload
}

// PhotoUpload хранит имя файла и байты загружаемой фотографии.
type PhotoUpload struct {
	Name    string
	Content []byte
}
