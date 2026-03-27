package requests

import (
	"grovia/internal/models"
	"mime/multipart"
	"strings"
)

type LocationRequest struct {
	Name    string                `form:"name" validate:"required"`
	Address string                `form:"address" validate:"required"`
	Picture *multipart.FileHeader `form:"picture"`
}

func (r *LocationRequest) ToModel(pictureURL string) models.Location {
	location := models.Location{
		Name:    strings.TrimSpace(r.Name),
		Address: strings.TrimSpace(r.Address),
	}

	if pictureURL != "" {
		location.Picture = pictureURL
	}

	return location
}
