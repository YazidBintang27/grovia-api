package responses

import (
	"grovia/internal/models"
	"time"
)

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

func FromModelUser(user models.User) *UserResponse {
	return &UserResponse{
		ID:             user.ID,
		LocationID:     user.LocationID,
		Name:           user.Name,
		PhoneNumber:    user.PhoneNumber,
		Address:        user.Address,
		Nik:            user.Nik,
		ProfilePicture: user.ProfilePicture,
		Role:           user.Role,
		IsActive:       user.IsActive,
		CreatedBy:      user.CreatedBy,
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
	}
}

func FromModelUserList(users []models.User) []UserResponse {
	var responses []UserResponse
	for _, v := range users {
		responses = append(responses, *FromModelUser(v))
	}
	return responses
}
