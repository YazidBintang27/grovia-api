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

	toddlerMapping := models.Toddler{
		UpdatedByID: userID,
	}
	if req.Name != nil {
		toddlerMapping.Name = *req.Name
	}
	if req.PhoneNumber != nil {
		toddlerMapping.ParentID = parentID
	}
	if url != "" {
		toddlerMapping.ProfilePicture = url
	}
	if req.LocationID != nil {
		toddlerMapping.LocationID = *req.LocationID
	}
	toddler, err := t.repo.UpdateToddlerByID(id, locationID, &toddlerMapping)
	if err != nil {
		return nil, pkg.NewInternalServerError("Gagal update data toddler")
	}

	toddlerResponse := responses.ToddlerResponse{
		ID:                toddler.ID,
		ParentID:          toddler.ParentID,
		LocationID:        toddler.LocationID,
		CreatedByID:       toddler.CreatedByID,
		UpdatedByID:       toddler.UpdatedByID,
		Name:              toddler.Name,
		Birthdate:         toddler.Birthdate,
		Sex:               toddler.Sex,
		Height:            toddler.Height,
		NutritionalStatus: toddler.NutritionalStatus,
		ProfilePicture:    toddler.ProfilePicture,
		CreatedAt:         toddler.CreatedAt,
		UpdatedAt:         toddler.UpdatedAt,
	}

	return &toddlerResponse, nil
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

	var toddlerResponses []responses.ToddlerResponse
	for _, v := range toddlers {
		toddlerResponses = append(toddlerResponses, responses.ToddlerResponse{
			ID:                v.ID,
			ParentID:          v.ParentID,
			LocationID:        v.LocationID,
			CreatedByID:       v.CreatedByID,
			UpdatedByID:       v.UpdatedByID,
			Name:              v.Name,
			Birthdate:         v.Birthdate,
			Sex:               v.Sex,
			Height:            v.Height,
			ProfilePicture:    v.ProfilePicture,
			NutritionalStatus: v.NutritionalStatus,
			CreatedAt:         v.CreatedAt,
			UpdatedAt:         v.UpdatedAt,
		})
	}

	meta := responses.PaginationMeta{
		Page:      page,
		Limit:     limit,
		TotalData: total,
		TotalPage: totalPage,
	}

	return toddlerResponses, &meta, nil
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

	toddlerMapping := models.Toddler{
		ParentID:    parent.ID,
		CreatedByID: userID,
		UpdatedByID: userID,
		DeletedByID: nil,
		Name:        req.Name,
		Birthdate:   req.Birthdate,
		Sex:         req.Sex,
		Height:      req.Height,
		LocationID:  parent.LocationID,
	}

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

	toddlerModel := models.Toddler{
		ParentID:          parent.ID,
		LocationID:        parent.LocationID,
		Name:              toddler.Name,
		Birthdate:         toddler.Birthdate,
		Height:            toddler.Height,
		Sex:               toddler.Sex,
		ProfilePicture:    toddler.ProfilePicture,
		NutritionalStatus: predict.NutritionalStatus,
	}

	_, err = t.repo.UpdateToddlerByID(toddler.ID, parent.LocationID, &toddlerModel)

	if err != nil {
		return nil, nil, pkg.NewInternalServerError("Gagal update nutritional status")
	}

	toddlerResponse := responses.ToddlerResponse{
		ID:                toddler.ID,
		ParentID:          toddler.ParentID,
		LocationID:        toddler.LocationID,
		CreatedByID:       toddler.CreatedByID,
		UpdatedByID:       toddler.UpdatedByID,
		Name:              toddler.Name,
		Birthdate:         toddler.Birthdate,
		Sex:               toddler.Sex,
		Height:            toddler.Height,
		NutritionalStatus: predict.NutritionalStatus,
		CreatedAt:         toddler.CreatedAt,
		UpdatedAt:         toddler.UpdatedAt,
	}

	predictResponse := responses.PredictResponse{
		ID:                predict.ID,
		ToddlerID:         predict.ToddlerID,
		CreatedByID:       predict.CreatedByID,
		Height:            predict.Height,
		Age:               predict.Age,
		Sex:               predict.Sex,
		Zscore:            predict.Zscore,
		NutritionalStatus: predict.NutritionalStatus,
		CreatedAt:         predict.CreatedAt,
		UpdatedAt:         predict.UpdatedAt,
	}

	return &toddlerResponse, &predictResponse, nil
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

	parentMapping := models.Parent{
		CreatedByID: userID,
		UpdatedByID: userID,
		DeletedByID: nil,
		Name:        parentReq.Name,
		PhoneNumber: parentReq.PhoneNumber,
		Address:     parentReq.Address,
		Nik:         parentReq.Nik,
		Job:         parentReq.Job,
		LocationID:  parentReq.LocationID,
	}

	parent, err := t.parentRepo.CreateParent(&parentMapping)
	if err != nil {
		return nil, nil, nil, pkg.NewInternalServerError("Gagal membuat parent")
	}

	toddlerMapping := models.Toddler{
		ParentID:    parent.ID,
		CreatedByID: userID,
		UpdatedByID: userID,
		DeletedByID: nil,
		Name:        toddlerReq.Name,
		Birthdate:   toddlerReq.Birthdate,
		Sex:         toddlerReq.Sex,
		Height:      toddlerReq.Height,
		LocationID:  toddlerReq.LocationID,
	}

	toddler, err := t.repo.CreateToddler(&toddlerMapping)
	if err != nil {
		return nil, nil, nil, pkg.NewInternalServerError("Gagal membuat toddler")
	}

	predict, err := t.predict.CreateIndividualPredict(toddlerReq, parentReq.LocationID, toddler.ID, userID)
	if err != nil {
		return nil, nil, nil, pkg.NewInternalServerError("Gagal membuat prediksi")
	}

	toddlerModel := models.Toddler{
		ParentID:          parent.ID,
		LocationID:        parent.LocationID,
		CreatedByID:       userID,
		UpdatedByID:       userID,
		DeletedByID:       nil,
		Name:              toddler.Name,
		Birthdate:         toddler.Birthdate,
		Height:            toddler.Height,
		Sex:               toddler.Sex,
		ProfilePicture:    toddler.ProfilePicture,
		NutritionalStatus: predict.NutritionalStatus,
	}

	_, err = t.repo.UpdateToddlerByID(toddler.ID, parentReq.LocationID, &toddlerModel)
	if err != nil {
		return nil, nil, nil, pkg.NewInternalServerError("Gagal update nutritional status")
	}

	toddler.NutritionalStatus = predict.NutritionalStatus

	toddlerResponse := responses.ToddlerResponse{
		ID:                toddler.ID,
		ParentID:          toddler.ParentID,
		LocationID:        toddler.LocationID,
		CreatedByID:       toddler.CreatedByID,
		UpdatedByID:       toddler.UpdatedByID,
		Name:              toddler.Name,
		Birthdate:         toddler.Birthdate,
		Sex:               toddler.Sex,
		Height:            toddler.Height,
		ProfilePicture:    toddler.ProfilePicture,
		NutritionalStatus: toddler.NutritionalStatus,
		CreatedAt:         toddler.CreatedAt,
		UpdatedAt:         toddler.UpdatedAt,
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

	predictResponse := responses.PredictResponse{
		ID:                predict.ID,
		ToddlerID:         predict.ToddlerID,
		CreatedByID:       predict.CreatedByID,
		Height:            predict.Height,
		Age:               predict.Age,
		Sex:               predict.Sex,
		Zscore:            predict.Zscore,
		NutritionalStatus: predict.NutritionalStatus,
		CreatedAt:         predict.CreatedAt,
		UpdatedAt:         predict.UpdatedAt,
	}

	return &toddlerResponse, &parentResp, &predictResponse, nil
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

	var toddlerResponse []responses.ToddlerResponse

	for _, v := range toddlers {
		toddlerResponse = append(toddlerResponse, responses.ToddlerResponse{
			ID:                v.ID,
			ParentID:          v.ParentID,
			LocationID:        v.LocationID,
			CreatedByID:       v.CreatedByID,
			UpdatedByID:       v.UpdatedByID,
			Name:              v.Name,
			Birthdate:         v.Birthdate,
			Sex:               v.Sex,
			Height:            v.Height,
			ProfilePicture:    v.ProfilePicture,
			NutritionalStatus: v.NutritionalStatus,
			CreatedAt:         v.CreatedAt,
			UpdatedAt:         v.UpdatedAt,
		})
	}

	meta := responses.PaginationMeta{
		Page:      page,
		Limit:     limit,
		TotalData: total,
		TotalPage: totalPage,
	}

	return toddlerResponse, &meta, err
}

func (t *toddlerService) GetToddlerByID(id int, locationID int) (*responses.ToddlerResponse, error) {
	toddler, err := t.repo.GetToddlerByID(id, locationID)

	if err != nil {
		return nil, pkg.NewNotFoundError("Toddler tidak ditemukan")
	}

	toddlerResponse := responses.ToddlerResponse{
		ID:                toddler.ID,
		ParentID:          toddler.ParentID,
		LocationID:        toddler.LocationID,
		CreatedByID:       toddler.CreatedByID,
		UpdatedByID:       toddler.UpdatedByID,
		Name:              toddler.Name,
		Birthdate:         toddler.Birthdate,
		Sex:               toddler.Sex,
		Height:            toddler.Height,
		ProfilePicture:    toddler.ProfilePicture,
		NutritionalStatus: toddler.NutritionalStatus,
		CreatedAt:         toddler.CreatedAt,
		UpdatedAt:         toddler.UpdatedAt,
	}

	return &toddlerResponse, nil
}

func (t *toddlerService) UpdateToddlerByID(
	ctx context.Context, id, locationID, userID int,
	req requests.UpdateToddlerRequest,
) (*responses.ToddlerResponse, *responses.PredictResponse, error) {

	if err := pkg.ValidateStruct(req); err != nil {
		return nil, nil, pkg.NewBadRequestError(err.Error())
	}

	toddlerRequest := requests.CreateToddlerRequest{}
	if req.Name != nil {
		toddlerRequest.Name = *req.Name
	}
	if req.Birthdate != nil {
		toddlerRequest.Birthdate = *req.Birthdate
	}
	if req.Sex != "" {
		toddlerRequest.Sex = req.Sex
	}
	if req.Height != nil {
		toddlerRequest.Height = *req.Height
	}
	if req.NutritionalStatus != nil {
		toddlerRequest.NutritionalStatus = *req.NutritionalStatus
	}
	if req.LocationID != nil {
		toddlerRequest.LocationID = *req.LocationID
	}
	if req.PhoneNumber != nil {
		toddlerRequest.PhoneNumber = *req.PhoneNumber
	}

	predict, err := t.predict.CreateIndividualPredict(
		toddlerRequest,
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

	toddlerMapping := models.Toddler{
		UpdatedByID: userID,
	}
	if req.Name != nil {
		toddlerMapping.Name = *req.Name
	}
	if req.Birthdate != nil {
		toddlerMapping.Birthdate = *req.Birthdate
	}
	if req.Sex != "" {
		toddlerMapping.Sex = req.Sex
	}
	if req.Height != nil {
		toddlerMapping.Height = *req.Height
	}
	if req.PhoneNumber != nil {
		toddlerMapping.ParentID = parentID
	}
	if url != "" {
		toddlerMapping.ProfilePicture = url
	}
	if req.LocationID != nil {
		toddlerMapping.LocationID = *req.LocationID
	}
	toddlerMapping.NutritionalStatus = predict.NutritionalStatus

	toddler, err := t.repo.UpdateToddlerByID(id, locationID, &toddlerMapping)
	if err != nil {
		return nil, nil, pkg.NewInternalServerError("Gagal update data toddler")
	}

	toddlerResponse := responses.ToddlerResponse{
		ID:                toddler.ID,
		ParentID:          toddler.ParentID,
		LocationID:        toddler.LocationID,
		CreatedByID:       toddler.CreatedByID,
		UpdatedByID:       toddler.UpdatedByID,
		Name:              toddler.Name,
		Birthdate:         toddler.Birthdate,
		Sex:               toddler.Sex,
		Height:            toddler.Height,
		NutritionalStatus: toddler.NutritionalStatus,
		ProfilePicture:    toddler.ProfilePicture,
		CreatedAt:         toddler.CreatedAt,
		UpdatedAt:         toddler.UpdatedAt,
	}

	predictResponse := responses.PredictResponse{
		ID:                predict.ID,
		ToddlerID:         predict.ToddlerID,
		CreatedByID:       predict.CreatedByID,
		Name:              predict.Name,
		Height:            predict.Height,
		Age:               predict.Age,
		Sex:               predict.Sex,
		Zscore:            predict.Zscore,
		NutritionalStatus: predict.NutritionalStatus,
		CreatedAt:         predict.CreatedAt,
		UpdatedAt:         predict.UpdatedAt,
	}

	return &toddlerResponse, &predictResponse, nil
}

func NewToddlerService(repo repositories.ToddlerRepository, parentRepo repositories.ParentRepository, s3 *S3Service, predict PredictService) ToddlerService {
	return &toddlerService{repo: repo, parentRepo: parentRepo, s3: s3, predict: predict}
}
