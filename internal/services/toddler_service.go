package services

import (
	"context"
	"grovia/internal/dto/requests"
	"grovia/internal/dto/responses"
	"grovia/internal/models"
	"grovia/internal/repositories"
	"grovia/pkg"
	"math"
	"strconv"
)

type ToddlerService interface {
	CreateToddler(req requests.CreateToddlerRequest, userID int) (*responses.ToddlerResponse, *responses.PredictResponse, error)
	CreateToddlerWithParent(toddlerReq requests.CreateToddlerRequest, parentReq requests.CreateParentRequest, userID int) (*responses.ToddlerResponse, *responses.ParentResponse, *responses.PredictResponse, error)
	GetAllToddler(locationID int, name, pageStr, limitStr string) ([]responses.ToddlerResponse, *responses.PaginationMeta, error)
	GetToddlerByID(id, locationID int) (*responses.ToddlerResponse, error)
	UpdateToddlerByID(ctx context.Context, id, locationID, userID int, req requests.UpdateToddlerRequest) (*responses.ToddlerResponse, *responses.PredictResponse, error)
	DeleteToddlerByID(id, locationID, userID int) error
	CheckToddlerExists(phoneNumber, name string) (bool, *models.Toddler, error)
	GetAllToddlerAllLocation(name, pageStr, limitStr string) ([]responses.ToddlerResponse, *responses.PaginationMeta, error)
	UpdateToddlerByIDWithoutPredict(ctx context.Context, id, locationID, userID int, req requests.UpdateToddlerRequest) (*responses.ToddlerResponse, error)
}

type toddlerService struct {
	repo       repositories.ToddlerRepository
	parentRepo repositories.ParentRepository
	s3         *S3Service
	predict    PredictService
}

func (t *toddlerService) UpdateToddlerByIDWithoutPredict(ctx context.Context, id int, locationID, userID int, req requests.UpdateToddlerRequest) (*responses.ToddlerResponse, error) {
	if err := pkg.ValidateStruct(req); err != nil {
		return nil, pkg.NewBadRequestError(err.Error())
	}

	var url string
	var err error
	if req.ProfilePicture != nil && req.ProfilePicture.Filename != "" && req.ProfilePicture.Size > 0 {
		url, err = t.s3.UploadFile(ctx, req.ProfilePicture, "toddlers")
		if err != nil {
			return nil, pkg.NewInternalServerError("Gagal upload foto")
		}
	}

	var parentID int
	if req.PhoneNumber != nil {
		parent, err := t.parentRepo.FindParentByPhoneNumber(*req.PhoneNumber)
		if err != nil {
			return nil, pkg.NewNotFoundError("Orang tua dengan nomor HP " + *req.PhoneNumber + " tidak ditemukan")
		}
		parentID = parent.ID
	}

	toddlerMapping := req.ToUpdateModel(userID, url, parentID)

	toddler, err := t.repo.UpdateToddlerByID(id, locationID, &toddlerMapping)
	if err != nil {
		return nil, pkg.NewInternalServerError("Gagal update data toddler")
	}

	return responses.FromModelToddler(*toddler), nil
}

func (t *toddlerService) GetAllToddlerAllLocation(name, pageStr, limitStr string) ([]responses.ToddlerResponse, *responses.PaginationMeta, error) {
	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 1
	}

	offset := (page - 1) * limit

	toddlers, total, err := t.repo.GetAllToddlerAllLocation(name, limit, offset)

	if err != nil {
		return nil, nil, pkg.NewInternalServerError("Gagal mengambil data toddler")
	}

	totalPage := int(math.Ceil(float64(total) / float64(limit)))

	meta := responses.PaginationMeta{
		Page:      page,
		Limit:     limit,
		TotalData: total,
		TotalPage: totalPage,
	}

	return responses.FromModelToddlerList(toddlers), &meta, nil
}

func (t *toddlerService) CheckToddlerExists(phoneNumber, name string) (bool, *models.Toddler, error) {
	parent, err := t.parentRepo.FindParentByPhoneNumber(phoneNumber)
	if err != nil {
		return false, nil, pkg.NewNotFoundError("Parent tidak ditemukan")
	}
	return t.repo.FindToddlerByName(parent.ID, name)
}

func (t *toddlerService) CreateToddler(req requests.CreateToddlerRequest, userID int) (*responses.ToddlerResponse, *responses.PredictResponse, error) {
	parent, err := t.parentRepo.FindParentByPhoneNumber(req.PhoneNumber)

	if parent == nil {
		return nil, nil, pkg.NewNotFoundError("Parent dengan nomor telepon " + req.PhoneNumber + " tidak ditemukan")
	}

	toddlerMapping := req.ToModel(userID, parent.ID)

	if err != nil {
		return nil, nil, pkg.NewInternalServerError("Gagal memproses data parent")
	}

	toddler, err := t.repo.CreateToddler(&toddlerMapping)

	if err != nil {
		return nil, nil, pkg.NewInternalServerError("Gagal membuat toddler")
	}

	predict, err := t.predict.CreateIndividualPredict(req, toddler.LocationID, toddler.ID, userID)

	if err != nil {
		return nil, nil, pkg.NewInternalServerError("Gagal membuat prediksi")
	}

	toddler.LocationID = parent.LocationID
	toddler.NutritionalStatus = predict.NutritionalStatus
	toddlerModel := req.ToModel(userID, parent.ID)

	_, err = t.repo.UpdateToddlerByID(toddler.ID, parent.LocationID, &toddlerModel)

	if err != nil {
		return nil, nil, pkg.NewInternalServerError("Gagal update nutritional status")
	}

	toddler.NutritionalStatus = predict.NutritionalStatus

	return responses.FromModelToddler(*toddler), predict, nil
}

func (t *toddlerService) CreateToddlerWithParent(
	toddlerReq requests.CreateToddlerRequest,
	parentReq requests.CreateParentRequest,
	userID int,
) (*responses.ToddlerResponse, *responses.ParentResponse, *responses.PredictResponse, error) {
	if err := pkg.ValidateStruct(toddlerReq); err != nil {
		return nil, nil, nil, pkg.NewBadRequestError(err.Error())
	}

	if err := pkg.ValidateStruct(parentReq); err != nil {
		return nil, nil, nil, pkg.NewBadRequestError(err.Error())
	}

	parentMapping := parentReq.ToModel(userID)
	parent, err := t.parentRepo.CreateParent(&parentMapping)
	if err != nil {
		return nil, nil, nil, pkg.NewInternalServerError("Gagal membuat parent")
	}

	toddlerMapping := toddlerReq.ToModel(userID, parent.ID)

	toddler, err := t.repo.CreateToddler(&toddlerMapping)
	if err != nil {
		return nil, nil, nil, pkg.NewInternalServerError("Gagal membuat toddler")
	}

	predict, err := t.predict.CreateIndividualPredict(toddlerReq, parentReq.LocationID, toddler.ID, userID)
	if err != nil {
		return nil, nil, nil, pkg.NewInternalServerError("Gagal membuat prediksi")
	}

	toddler.LocationID = parent.LocationID
	toddler.NutritionalStatus = predict.NutritionalStatus

	toddlerModel := toddlerReq.ToModel(userID, parent.ID)

	_, err = t.repo.UpdateToddlerByID(toddler.ID, parentReq.LocationID, &toddlerModel)
	if err != nil {
		return nil, nil, nil, pkg.NewInternalServerError("Gagal update nutritional status")
	}

	toddler.NutritionalStatus = predict.NutritionalStatus

	return responses.FromModelToddler(*toddler), responses.FromModelParent(*parent), predict, nil
}

func (t *toddlerService) DeleteToddlerByID(id int, locationID, userID int) error {
	err := t.repo.DeleteToddlerByID(id, locationID, userID)
	if err != nil {
		return pkg.NewInternalServerError("Gagal menghapus toddler")
	}
	return nil
}

func (t *toddlerService) GetAllToddler(locationID int, name, pageStr, limitStr string) ([]responses.ToddlerResponse, *responses.PaginationMeta, error) {
	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 1
	}

	offset := (page - 1) * limit

	toddlers, total, err := t.repo.GetAllToddler(locationID, limit, offset, name)

	if err != nil {
		return nil, nil, pkg.NewInternalServerError("Gagal mengambil data toddler")
	}

	totalPage := int(math.Ceil(float64(total) / float64(limit)))

	meta := responses.PaginationMeta{
		Page:      page,
		Limit:     limit,
		TotalData: total,
		TotalPage: totalPage,
	}

	return responses.FromModelToddlerList(toddlers), &meta, err
}

func (t *toddlerService) GetToddlerByID(id int, locationID int) (*responses.ToddlerResponse, error) {
	toddler, err := t.repo.GetToddlerByID(id, locationID)

	if err != nil {
		return nil, pkg.NewNotFoundError("Toddler tidak ditemukan")
	}

	return responses.FromModelToddler(*toddler), nil
}

func (t *toddlerService) UpdateToddlerByID(
	ctx context.Context, id, locationID, userID int,
	req requests.UpdateToddlerRequest,
) (*responses.ToddlerResponse, *responses.PredictResponse, error) {

	if err := pkg.ValidateStruct(req); err != nil {
		return nil, nil, pkg.NewBadRequestError(err.Error())
	}

	predict, err := t.predict.CreateIndividualPredict(
		req.ToCreateRequest(),
		locationID,
		id,
		userID,
	)
	if err != nil {
		return nil, nil, pkg.NewInternalServerError("Gagal membuat prediksi")
	}

	var url string

	if req.ProfilePicture != nil && req.ProfilePicture.Filename != "" && req.ProfilePicture.Size > 0 {
		url, err = t.s3.UploadFile(ctx, req.ProfilePicture, "toddlers")
		if err != nil {
			return nil, nil, pkg.NewInternalServerError("Gagal upload foto")
		}
	}

	var parentID int
	if req.PhoneNumber != nil {
		parent, err := t.parentRepo.FindParentByPhoneNumber(*req.PhoneNumber)
		if err != nil {
			return nil, nil, pkg.NewNotFoundError("Orang tua dengan nomor HP " + *req.PhoneNumber + " tidak ditemukan")
		}
		parentID = parent.ID
	}

	toddlerMapping := req.ToUpdateModel(userID, url, parentID)
	toddlerMapping.NutritionalStatus = predict.NutritionalStatus

	toddler, err := t.repo.UpdateToddlerByID(id, locationID, &toddlerMapping)
	if err != nil {
		return nil, nil, pkg.NewInternalServerError("Gagal update data toddler")
	}

	return responses.FromModelToddler(*toddler), predict, nil
}

func NewToddlerService(repo repositories.ToddlerRepository, parentRepo repositories.ParentRepository, s3 *S3Service, predict PredictService) ToddlerService {
	return &toddlerService{repo: repo, parentRepo: parentRepo, s3: s3, predict: predict}
}
