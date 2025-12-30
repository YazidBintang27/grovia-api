package responses

import "time"

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
