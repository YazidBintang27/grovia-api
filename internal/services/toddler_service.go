package services

import (
	"grovia/internal/dto/requests"
	"grovia/internal/dto/responses"
	"grovia/internal/models"
	"grovia/internal/repositories"
)

type ToddlerService interface {
	GetAllToddler(locationID int) ([]responses.ToddlerResponse, error)
	GetToddlerByID(id, locationID int) (*responses.ToddlerResponse, error)
	UpdateToddlerByID(id, locationID int, req requests.UpdateToddlerRequest) (*responses.ToddlerResponse, error)
	DeleteToddlerByID(id, locationID int) error
}

type toddlerService struct {
	repo repositories.ToddlerRepository
	s3   *S3Service
}

// DeleteToddlerByID implements ToddlerService.
func (t *toddlerService) DeleteToddlerByID(id int, locationID int) error {
	return t.repo.DeleteToddlerByID(id, locationID)
}

// GetAllToddler implements ToddlerService.
func (t *toddlerService) GetAllToddler(locationID int) ([]responses.ToddlerResponse, error) {
	toddlers, err := t.repo.GetAllToddler(locationID)

	if err != nil {
		return nil, err
	}

	var toddlerResponse []responses.ToddlerResponse

	for _, v := range toddlers {
		toddlerResponse = append(toddlerResponse, responses.ToddlerResponse{
			ID:             v.ID,
			ParentID:       v.ParentID,
			LocationID:     v.Parent.LocationID,
			Name:           v.Name,
			Birthdate:      v.Birthdate,
			Gender:         v.Gender,
			Height:         v.Height,
			ProfilePicture: v.ProfilePicture,
			CreatedAt:      v.CreatedAt,
			UpdatedAt:      v.UpdatedAt,
		})
	}

	return toddlerResponse, err
}

// GetToddlerByID implements ToddlerService.
func (t *toddlerService) GetToddlerByID(id int, locationID int) (*responses.ToddlerResponse, error) {
	toddler, err := t.repo.GetToddlerByID(id, locationID)

	if err != nil {
		return nil, err
	}

	toddlerResponse := responses.ToddlerResponse{
		ID:             toddler.ID,
		ParentID:       toddler.ParentID,
		LocationID:     toddler.LocationID,
		Name:           toddler.Name,
		Birthdate:      toddler.Birthdate,
		Gender:         toddler.Gender,
		Height:         toddler.Height,
		ProfilePicture: toddler.ProfilePicture,
		CreatedAt:      toddler.CreatedAt,
		UpdatedAt:      toddler.UpdatedAt,
	}

	return &toddlerResponse, nil
}

// UpdateToddlerByID implements ToddlerService.
func (t *toddlerService) UpdateToddlerByID(id int, locationID int, req requests.UpdateToddlerRequest) (*responses.ToddlerResponse, error) {

	var err error
	var url string
	if req.ProfilePicture != nil {
		url, err = t.s3.UploadFile(req.ProfilePicture, "toddlers")
		if err != nil {
			return nil, err
		}
	}

	toddlerMapping := models.Toddler{
		Name:      req.Name,
		Birthdate: req.Birthdate,
		Gender:    req.Gender,
		Height:    req.Height,
	}

	if url != "" {
		toddlerMapping.ProfilePicture = url
	}

	toddler, err := t.repo.UpdateToddlerByID(id, locationID, &toddlerMapping)

	if err != nil {
		return nil, err
	}

	toddlerResponse := responses.ToddlerResponse{
		ID:             toddler.ID,
		ParentID:       toddler.ParentID,
		LocationID:     toddler.LocationID,
		Name:           toddler.Name,
		Birthdate:      toddler.Birthdate,
		Gender:         toddler.Gender,
		Height:         toddler.Height,
		ProfilePicture: toddler.ProfilePicture,
		CreatedAt:      toddler.CreatedAt,
		UpdatedAt:      toddler.UpdatedAt,
	}

	return &toddlerResponse, nil
}

func NewToddlerService(repo repositories.ToddlerRepository) ToddlerService {
	return &toddlerService{repo: repo}
}
