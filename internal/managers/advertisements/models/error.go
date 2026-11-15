package models

// ErrorResponse описывает структурированную ошибку API.
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}
