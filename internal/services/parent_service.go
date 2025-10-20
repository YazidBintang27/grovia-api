package services

import (
	"errors"
	"fmt"
	"grovia/internal/dto/requests"
	"grovia/internal/dto/responses"
	"grovia/internal/models"
	"grovia/internal/repositories"
	"regexp"
	"strings"
)

type ParentService interface {
	GetAllParent(locationID int) ([]responses.ParentResponse, error)
	GetParentByID(id, locationID int) (*responses.ParentResponse, error)
	UpdateParentByID(id, locationID int, req requests.UpdateParentRequest) (*responses.ParentResponse, error)
	DeleteParentByID(id, locationID int) error
	CheckPhoneExists(phoneNumber string) (*models.Parent, error)
	GetAllParentAllLocation() ([]responses.ParentResponse, error)
}

type parentService struct {
	repo repositories.ParentRepository
}

// GetAllParentAllLocation implements ParentService.
func (p *parentService) GetAllParentAllLocation() ([]responses.ParentResponse, error) {
	parents, err := p.repo.GetAllParentAllLocation()

	if err != nil {
		return nil, err
	}

	var parentResponses []responses.ParentResponse
	for _, v := range parents {
		parentResponses = append(parentResponses, responses.ParentResponse{
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

	return parentResponses, nil
}

// CheckPhoneExists implements ParentService.
func (p *parentService) CheckPhoneExists(phoneNumber string) (*models.Parent, error) {
	return p.repo.FindParentByPhoneNumber(phoneNumber)
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
	parentMapping := models.Parent{}

	if req.Name != nil {
		trimmed := strings.TrimSpace(*req.Name)
		parentMapping.Name = trimmed
	}

	if req.PhoneNumber != nil {
		phone := strings.TrimSpace(*req.PhoneNumber)
		if phone != "" {
			if len(phone) < 10 || len(phone) > 15 {
				return nil, errors.New("nomor HP harus antara 10–15 digit")
			}

			if !regexp.MustCompile(`^[0-9]+$`).MatchString(phone) {
				return nil, errors.New("nomor HP hanya boleh berisi angka")
			}
		}
		parentMapping.PhoneNumber = phone
	}

	if req.Address != nil {
		trimmed := strings.TrimSpace(*req.Address)
		parentMapping.Address = trimmed
	}

	if req.Nik != nil {
		nik := strings.TrimSpace(*req.Nik)
		if nik != "" {
			if !regexp.MustCompile(`^[0-9]{16}$`).MatchString(nik) {
				return nil, errors.New("NIK harus terdiri dari 16 digit angka")
			}
		}
		parentMapping.Nik = nik
	}

	if req.Job != nil {
		trimmed := strings.TrimSpace(*req.Job)
		parentMapping.Job = trimmed
	}

	if req.LocationID != nil {
		parentMapping.LocationID = *req.LocationID
	}

	parent, err := p.repo.UpdateParentByID(id, locationID, &parentMapping)
	if err != nil {
		return nil, fmt.Errorf("gagal update data orang tua: %w", err)
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
