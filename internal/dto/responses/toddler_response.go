package responses

import (
	"grovia/internal/models"
	"time"
)

type ToddlerResponse struct {
	ID                int       `json:"id"`
	ParentID          int       `json:"parentID"`
	LocationID        int       `json:"locationID"`
	CreatedByID       int       `json:"createdByID"`
	UpdatedByID       int       `json:"updatedByID"`
	Name              string    `json:"name"`
	Birthdate         time.Time `json:"birthdate"`
	Sex               string    `json:"sex"`
	Height            float64   `json:"height"`
	ProfilePicture    string    `json:"profilePicture"`
	NutritionalStatus string    `json:"nutritionalStatus"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}

func FromModelToddler(toddler models.Toddler) *ToddlerResponse {
	return &ToddlerResponse{
		ID:                toddler.ID,
		ParentID:          toddler.ParentID,
		LocationID:        toddler.LocationID,
		CreatedByID:       toddler.CreatedByID,
		UpdatedByID:       toddler.UpdatedByID,
		Name:              toddler.Name,
		Birthdate:         toddler.Birthdate,
		Sex:               toddler.Sex,
		Height:            toddler.Height,
		NutritionalStatus: toddler.NutritionalStatus,
		CreatedAt:         toddler.CreatedAt,
		UpdatedAt:         toddler.UpdatedAt,
	}
}

func FromModelToddlerList(toddlers []models.Toddler) []ToddlerResponse {
	var toddlerResponse []ToddlerResponse
	for _, t := range toddlers {
		toddlerResponse = append(toddlerResponse, *FromModelToddler(t))
	}
	return toddlerResponse
}
