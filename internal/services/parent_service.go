package services

import (
	"grovia/internal/dto/requests"
	"grovia/internal/dto/responses"
	"grovia/internal/models"
	"grovia/internal/repositories"
)

type ParentService interface {
	CreateParentWithToddlers(req requests.CreateParentWithToddlersRequest) (*responses.ParentWithToddlersResponse, error)
	GetAllParent(locationID int) ([]responses.ParentResponse, error)
	GetParentByID(id, locationID int) (*responses.ParentResponse, error)
	UpdateParentByID(id, locationID int, req requests.UpdateParentRequest) (*responses.ParentResponse, error)
	DeleteParentByID(id, locationID int) error
}

type parentService struct {
	repo repositories.ParentRepository
}

// CreateParentWithToddlers implements ParentService.
func (p *parentService) CreateParentWithToddlers(req requests.CreateParentWithToddlersRequest) (*responses.ParentWithToddlersResponse, error) {
	parentMapping := models.Parent{
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
		Address:     req.Address,
		Nik:         req.Nik,
		Job:         req.Job,
	}

	if req.LocationID != nil {
		parentMapping.LocationID = *req.LocationID
	}

	var toddlers []models.Toddler
	for _, t := range req.Toddlers {
		toddlers = append(toddlers, models.Toddler{
			Name:      t.Name,
			Birthdate: t.Birthdate,
			Gender:    t.Gender,
			Height:    t.Height,
		})
	}

	parent, err := p.repo.CreateParentWithToddlers(&parentMapping, toddlers)
	if err != nil {
		return nil, err
	}

	var toddlerResponses []responses.ToddlerResponse
	for _, t := range toddlers {
		toddlerResponses = append(toddlerResponses, responses.ToddlerResponse{
			ID:             t.ID,
			ParentID:       t.ParentID,
			Name:           t.Name,
			Birthdate:      t.Birthdate,
			Gender:         t.Gender,
			Height:         t.Height,
			ProfilePicture: t.ProfilePicture,
			CreatedAt:      t.CreatedAt,
			UpdatedAt:      t.UpdatedAt,
		})
	}

	parentResponse := responses.ParentWithToddlersResponse{
		ID:          parent.ID,
		LocationID:  parent.LocationID,
		Name:        parent.Name,
		PhoneNumber: parent.PhoneNumber,
		Address:     parent.Address,
		Nik:         parent.Nik,
		Job:         parent.Job,
		Toddlers:    toddlerResponses,
		CreatedAt:   parent.CreatedAt,
		UpdatedAt:   parent.UpdatedAt,
	}

	return &parentResponse, nil
}

// DeleteParentByID implements ParentService.
func (p *parentService) DeleteParentByID(id int, locationID int) error {
	return p.repo.DeleteParentByID(id, locationID)
}

// GetAllParent implements ParentService.
func (p *parentService) GetAllParent(locationID int) ([]responses.ParentResponse, error) {
	parents, err := p.repo.GetAllParent(locationID)

	if err != nil {
		return nil, err
	}

	var parentResponse []responses.ParentResponse

	for _, v := range parents {
		parentResponse = append(parentResponse, responses.ParentResponse{
			ID:          v.ID,
			LocationID:  v.LocationID,
			Name:        v.Name,
			PhoneNumber: v.PhoneNumber,
			Address:     v.Address,
			Nik:         v.Nik,
			Job:         v.Job,
			CreatedAt:   v.CreatedAt,
			UpdatedAt:   v.UpdatedAt,
		})
	}

	return parentResponse, nil
}

// GetParentByID implements ParentService.
func (p *parentService) GetParentByID(id int, locationID int) (*responses.ParentResponse, error) {
	parent, err := p.repo.GetParentByID(id, locationID)

	if err != nil {
		return nil, err
	}

	parentResponses := responses.ParentResponse{
		ID:          parent.ID,
		LocationID:  parent.LocationID,
		Name:        parent.Name,
		PhoneNumber: parent.PhoneNumber,
		Address:     parent.Address,
		Nik:         parent.Nik,
		Job:         parent.Job,
		CreatedAt:   parent.CreatedAt,
		UpdatedAt:   parent.UpdatedAt,
	}

	return &parentResponses, nil
}

// UpdateParentByID implements ParentService.
func (p *parentService) UpdateParentByID(id int, locationID int, req requests.UpdateParentRequest) (*responses.ParentResponse, error) {
	parentMapping := models.Parent{
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
		Address:     req.Address,
		Nik:         req.Nik,
		Job:         req.Job,
	}

	parent, err := p.repo.UpdateParentByID(id, locationID, &parentMapping)

	if err != nil {
		return nil, err
	}

	parentResponse := responses.ParentResponse{
		ID:          parent.ID,
		LocationID:  parent.LocationID,
		Name:        parent.Name,
		PhoneNumber: parent.PhoneNumber,
		Address:     parent.Address,
		Nik:         parent.Nik,
		Job:         parent.Job,
		CreatedAt:   parent.CreatedAt,
		UpdatedAt:   parent.UpdatedAt,
	}

	return &parentResponse, nil
}

func NewParentService(repo repositories.ParentRepository) ParentService {
	return &parentService{repo: repo}
}
