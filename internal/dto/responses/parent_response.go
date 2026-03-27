package responses

import (
	"grovia/internal/models"
	"time"
)

type ParentResponse struct {
	ID          int               `json:"id"`
	LocationID  int               `json:"locationID"`
	CreatedByID int               `json:"createdByID"`
	UpdatedByID int               `json:"updatedByID"`
	Name        string            `json:"name"`
	PhoneNumber string            `json:"phoneNumber"`
	Address     string            `json:"address"`
	Nik         string            `json:"nik"`
	Job         string            `json:"job"`
	Toddlers    []ToddlerResponse `json:"toddlers"`
	CreatedAt   time.Time         `json:"createdAt"`
	UpdatedAt   time.Time         `json:"updatedAt"`
}

func FromModelParent(parent models.Parent) *ParentResponse {
	toddlerResponses := FromModelToddlerList(parent.Toddlers)
	return &ParentResponse{
		ID:          parent.ID,
		LocationID:  parent.LocationID,
		CreatedByID: parent.CreatedByID,
		UpdatedByID: parent.UpdatedByID,
		Name:        parent.Name,
		PhoneNumber: parent.PhoneNumber,
		Address:     parent.Address,
		Nik:         parent.Nik,
		Job:         parent.Job,
		Toddlers:    toddlerResponses,
		CreatedAt:   parent.CreatedAt,
		UpdatedAt:   parent.UpdatedAt,
	}
}

func FromModelParentList(parents []models.Parent) []ParentResponse {
	var parentResponse []ParentResponse
	for _, v := range parents {
		parentResponse = append(parentResponse, *FromModelParent(v))
	}
	return parentResponse
}
