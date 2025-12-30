package responses

import "time"

type UserResponse struct {
	ID             int       `json:"id"`
	LocationID     int       `json:"locationID"`
	Name           string    `json:"name"`
	PhoneNumber    string    `json:"phoneNumber"`
	Address        string    `json:"address"`
	Nik            string    `json:"nik"`
	ProfilePicture string    `json:"profilePicture"`
	Role           string    `json:"role"`
	IsActive       bool      `json:"isActive"`
	CreatedBy      string    `json:"createdBy"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}
