package models

import advertisementModels "api-tests-template/internal/managers/advertisements/models"

type Favorite struct {
	FavoriteID    string                            `json:"favorite_id"`
	Advertisement advertisementModels.Advertisement `json:"advertisement"`
}

type FavoritesList struct {
	Items []Favorite `json:"items"`
}
