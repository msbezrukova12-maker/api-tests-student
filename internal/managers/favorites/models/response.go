package models

import advertisementModels "api-tests-template/internal/managers/advertisements/models"

// Favorite содержит объявление из списка избранного.
type Favorite struct {
	FavoriteID    string                            `json:"favorite_id"`
	AddedAt       string                            `json:"added_at"`
	Advertisement advertisementModels.Advertisement `json:"advertisement"`
}

// FavoritesList описывает ответ GET /favorites.
type FavoritesList struct {
	Items []Favorite `json:"items"`
}
