package responses

import "time"

type ToddlerResponse struct {
	ID             int       `json:"id"`
	ParentID       int       `json:"parent_id"`
	LocationID     int       `json:"location_id"`
	Name           string    `json:"name"`
	Birthdate      time.Time `json:"birthdate"`
	Gender         string    `json:"gender"`
	Height         float64   `json:"height"`
	ProfilePicture string    `json:"profile_picture"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
