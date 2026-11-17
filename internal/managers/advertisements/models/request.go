package models

type CreateAdvertisementRequest struct {
	Fields map[string]string
	Photos []PhotoUpload
}

type PhotoUpload struct {
	Name    string
	Content []byte
}
