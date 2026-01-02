package services

import (
	"errors"
	"fmt"
	"grovia/internal/dto/requests"
	"grovia/internal/dto/responses"
	"grovia/internal/models"
	"grovia/internal/repositories"
	"math"
	"regexp"
	"strconv"
	"strings"
)

type ParentService interface {
	CreateParent(req requests.CreateParentRequest, userID int) (*responses.ParentResponse, error)
	GetAllParent(locationID int, name, pageStr, limitStr string) ([]responses.ParentResponse, *responses.PaginationMeta, error)
	GetParentByID(id, locationID int) (*responses.ParentResponse, error)
	UpdateParentByID(id, locationID, userID int, req requests.UpdateParentRequest) (*responses.ParentResponse, error)
	DeleteParentByID(id, locationID, userID int) error
	CheckPhoneExists(phoneNumber string) (*models.Parent, error)
	GetAllParentAllLocation(name, pageStr, limitStr string) ([]responses.ParentResponse, *responses.PaginationMeta, error)
}

type parentService struct {
	repo repositories.ParentRepository
}

// CreateParent implements ParentService.
func (p *parentService) CreateParent(req requests.CreateParentRequest, userID int) (*responses.ParentResponse, error) {
	if strings.TrimSpace(req.Name) == "" ||
		strings.TrimSpace(req.PhoneNumber) == "" ||
		strings.TrimSpace(req.Address) == "" ||
		strings.TrimSpace(req.Nik) == "" ||
		strings.TrimSpace(req.Job) == "" {
		return nil, fmt.Errorf("semua field parent (name, phone_number, address, nik, job) wajib diisi")
	}

	if len(req.PhoneNumber) < 10 || len(req.PhoneNumber) > 15 {
		return nil, fmt.Errorf("nomor telepon harus memiliki panjang antara 10 sampai 15 digit")
	}

	if len(req.Nik) != 16 {
		return nil, fmt.Errorf("NIK harus memiliki tepat 16 digit")
	}

	parentMapping := models.Parent{
		CreatedByID: userID,
		UpdatedByID: userID,
		DeletedByID: nil,
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
		Address:     req.Address,
		Nik:         req.Nik,
		Job:         req.Job,
		LocationID:  req.LocationID,
	}

	parent, err := p.repo.CreateParent(&parentMapping)

	if err != nil {
		return nil, err
	}

	parentResp := responses.ParentResponse{
		ID:          parent.ID,
		LocationID:  parent.LocationID,
		CreatedByID: parent.CreatedByID,
		UpdatedByID: parent.UpdatedByID,
		Name:        parent.Name,
		PhoneNumber: parent.PhoneNumber,
		Address:     parent.Address,
		Nik:         parent.Nik,
		Job:         parent.Job,
		CreatedAt:   parent.CreatedAt,
		UpdatedAt:   parent.UpdatedAt,
	}

	return &parentResp, nil
}

// GetAllParentAllLocation implements ParentService.
func (p *parentService) GetAllParentAllLocation(name, pageStr, limitStr string) ([]responses.ParentResponse, *responses.PaginationMeta, error) {
	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 1
	}

	offset := (page - 1) * limit

	parents, total, err := p.repo.GetAllParentAllLocation(name, limit, offset)

	if err != nil {
		return nil, nil, err
	}

	totalPage := int(math.Ceil(float64(total) / float64(limit)))

	var parentResponses []responses.ParentResponse
	for _, v := range parents {
		parentResponses = append(parentResponses, responses.ParentResponse{
			ID:          v.ID,
			LocationID:  v.LocationID,
			CreatedByID: v.CreatedByID,
			UpdatedByID: v.UpdatedByID,
			Name:        v.Name,
			PhoneNumber: v.PhoneNumber,
			Address:     v.Address,
			Nik:         v.Nik,
			Job:         v.Job,
			CreatedAt:   v.CreatedAt,
			UpdatedAt:   v.UpdatedAt,
		})
	}

	meta := responses.PaginationMeta{
		Page:      page,
		Limit:     limit,
		TotalData: total,
		TotalPage: totalPage,
	}

	return parentResponses, &meta, nil
}

// CheckPhoneExists implements ParentService.
func (p *parentService) CheckPhoneExists(phoneNumber string) (*models.Parent, error) {
	return p.repo.FindParentByPhoneNumber(phoneNumber)
}

// DeleteParentByID implements ParentService.
func (p *parentService) DeleteParentByID(id int, locationID, userID int) error {
	return p.repo.DeleteParentByID(id, locationID, userID)
}

// GetAllParent implements ParentService.
func (p *parentService) GetAllParent(locationID int, name, pageStr, limitStr string) ([]responses.ParentResponse, *responses.PaginationMeta, error) {
	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 1
	}

	offset := (page - 1) * limit

	parents, total, err := p.repo.GetAllParent(locationID, limit, offset, name)

	if err != nil {
		return nil, nil, err
	}

	totalPage := int(math.Ceil(float64(total) / float64(limit)))

	var parentResponse []responses.ParentResponse

	for _, v := range parents {
		parentResponse = append(parentResponse, responses.ParentResponse{
			ID:          v.ID,
			LocationID:  v.LocationID,
			CreatedByID: v.CreatedByID,
			UpdatedByID: v.UpdatedByID,
			Name:        v.Name,
			PhoneNumber: v.PhoneNumber,
			Address:     v.Address,
			Nik:         v.Nik,
			Job:         v.Job,
			CreatedAt:   v.CreatedAt,
			UpdatedAt:   v.UpdatedAt,
		})
	}

	meta := responses.PaginationMeta{
		Page:      page,
		Limit:     limit,
		TotalData: total,
		TotalPage: totalPage,
	}

	return parentResponse, &meta, nil
}

// GetParentByID implements ParentService.
func (p *parentService) GetParentByID(id int, locationID int) (*responses.ParentResponse, error) {
	parent, err := p.repo.GetParentByID(id, locationID)

	if err != nil {
		return nil, err
	}

	var toddlerResponses []responses.ToddlerResponse
	for _, t := range parent.Toddlers {
		toddlerResponses = append(toddlerResponses, responses.ToddlerResponse{
			ID:                t.ID,
			ParentID:          t.ParentID,
			LocationID:        t.LocationID,
			CreatedByID:       t.CreatedByID,
			UpdatedByID:       t.UpdatedByID,
			Name:              t.Name,
			Birthdate:         t.Birthdate,
			Sex:               t.Sex,
			Height:            t.Height,
			ProfilePicture:    t.ProfilePicture,
			NutritionalStatus: t.NutritionalStatus,
			CreatedAt:         t.CreatedAt,
			UpdatedAt:         t.UpdatedAt,
		})
	}

	parentResponses := responses.ParentResponse{
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

	return &parentResponses, nil
}

// UpdateParentByID implements ParentService.
func (p *parentService) UpdateParentByID(id int, locationID, userID int, req requests.UpdateParentRequest) (*responses.ParentResponse, error) {
	parentMapping := models.Parent{
		UpdatedByID: userID,
	}

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
		CreatedByID: parent.CreatedByID,
		UpdatedByID: parent.UpdatedByID,
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
