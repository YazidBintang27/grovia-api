package services

import (
	"grovia/internal/dto/requests"
	"grovia/internal/dto/responses"
	"grovia/internal/models"
	"grovia/internal/repositories"
	"grovia/pkg"
	"math"
	"strconv"
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

func (p *parentService) CreateParent(req requests.CreateParentRequest, userID int) (*responses.ParentResponse, error) {
	if err := pkg.ValidateStruct(req); err != nil {
		return nil, pkg.NewBadRequestError(err.Error())
	}

	parentMapping := req.ToModel(userID)

	parent, err := p.repo.CreateParent(&parentMapping)

	if err != nil {
		return nil, pkg.NewInternalServerError("Gagal membuat parent")
	}

	return responses.FromModelParent(*parent), nil
}

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
		return nil, nil, pkg.NewInternalServerError("Gagal mengambil data parent")
	}

	totalPage := int(math.Ceil(float64(total) / float64(limit)))

	meta := responses.PaginationMeta{
		Page:      page,
		Limit:     limit,
		TotalData: total,
		TotalPage: totalPage,
	}

	return responses.FromModelParentList(parents), &meta, nil
}

func (p *parentService) CheckPhoneExists(phoneNumber string) (*models.Parent, error) {
	parent, err := p.repo.FindParentByPhoneNumber(phoneNumber)
	if err != nil {
		return nil, pkg.NewNotFoundError("Nomor telepon tidak ditemukan")
	}
	return parent, nil
}

func (p *parentService) DeleteParentByID(id int, locationID, userID int) error {
	err := p.repo.DeleteParentByID(id, locationID, userID)
	if err != nil {
		return pkg.NewInternalServerError("Gagal menghapus parent")
	}
	return nil
}

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
		return nil, nil, pkg.NewInternalServerError("Gagal mengambil data parent")
	}

	totalPage := int(math.Ceil(float64(total) / float64(limit)))

	meta := responses.PaginationMeta{
		Page:      page,
		Limit:     limit,
		TotalData: total,
		TotalPage: totalPage,
	}

	return responses.FromModelParentList(parents), &meta, nil
}

func (p *parentService) GetParentByID(id int, locationID int) (*responses.ParentResponse, error) {
	parent, err := p.repo.GetParentByID(id, locationID)

	if err != nil {
		return nil, pkg.NewNotFoundError("Parent tidak ditemukan")
	}

	return responses.FromModelParent(*parent), nil
}

func (p *parentService) UpdateParentByID(id int, locationID, userID int, req requests.UpdateParentRequest) (*responses.ParentResponse, error) {
	if err := pkg.ValidateStruct(req); err != nil {
		return nil, pkg.NewBadRequestError(err.Error())
	}

	parentMapping := req.ToUpdateModel(userID)

	parent, err := p.repo.UpdateParentByID(id, locationID, &parentMapping)
	if err != nil {
		return nil, pkg.NewInternalServerError("Gagal update data parent")
	}

	return responses.FromModelParent(*parent), nil
}

func NewParentService(repo repositories.ParentRepository) ParentService {
	return &parentService{repo: repo}
}
