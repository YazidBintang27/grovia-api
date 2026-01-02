package requests

import "mime/multipart"

type LocationRequest struct {
	Name    string                `form:"name" validate:"required"`
	Address string                `form:"address" validate:"required"`
	Picture *multipart.FileHeader `form:"picture"`
}
