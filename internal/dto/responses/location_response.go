package responses

import (
	"grovia/internal/models"
	"time"
)

type LocationResponse struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	Picture   string    `json:"picture"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func FromModelLocation(location models.Location) *LocationResponse {
	return &LocationResponse{
		ID:        location.ID,
		Name:      location.Name,
		Address:   location.Address,
		Picture:   location.Picture,
		CreatedAt: location.CreatedAt,
		UpdatedAt: location.UpdatedAt,
	}
}

func FromModelLocationList(locations []models.Location) []LocationResponse {
	var locationResponse []LocationResponse
	for _, v := range locations {
		locationResponse = append(locationResponse, *FromModelLocation(v))
	}
	return locationResponse
}
