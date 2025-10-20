package services

import (
	"fmt"
	"grovia/internal/dto/requests"
	"grovia/internal/dto/responses"
	"grovia/internal/models"
	"grovia/internal/repositories"
	"strings"
)

type ToddlerService interface {
	CreateToddler(req requests.CreateToddlerRequest) (*responses.ToddlerResponse, *responses.PredictResponse, error)
	CreateToddlerWithParent(toddlerReq requests.CreateToddlerRequest, parentReq requests.CreateParentRequest) (*responses.ToddlerResponse, *responses.ParentResponse, *responses.PredictResponse, error)
	GetAllToddler(locationID int) ([]responses.ToddlerResponse, error)
	GetToddlerByID(id, locationID int) (*responses.ToddlerResponse, error)
	UpdateToddlerByID(id, locationID int, req requests.UpdateToddlerRequest) (*responses.ToddlerResponse, *responses.PredictResponse, error)
	DeleteToddlerByID(id, locationID int) error
	CheckToddlerExists(phoneNumber, name string) (bool, *models.Toddler, error)
	GetAllToddlerAllLocation() ([]responses.ToddlerResponse, error)
}

type toddlerService struct {
	repo       repositories.ToddlerRepository
	parentRepo repositories.ParentRepository
	s3         *S3Service
	predict    PredictService
}

// GetAllToddlerAllLocation implements ToddlerService.
func (t *toddlerService) GetAllToddlerAllLocation() ([]responses.ToddlerResponse, error) {
	toddlers, err := t.repo.GetAllToddlerAllLocation()

	if err != nil {
		return nil, err
	}

	var toddlerResponses []responses.ToddlerResponse
	for _, v := range toddlers {
		toddlerResponses = append(toddlerResponses, responses.ToddlerResponse{
			ID:                v.ID,
			ParentID:          v.ParentID,
			LocationID:        v.LocationID,
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

	return toddlerResponses, nil
}

// CheckToddlerExists implements ToddlerService.
func (t *toddlerService) CheckToddlerExists(phoneNumber, name string) (bool, *models.Toddler, error) {
	parent, err := t.parentRepo.FindParentByPhoneNumber(phoneNumber)
	if err != nil {
		return false, nil, err
	}
	return t.repo.FindToddlerByName(parent.ID, name)
}

// CreateToddler implements ToddlerService.
func (t *toddlerService) CreateToddler(req requests.CreateToddlerRequest) (*responses.ToddlerResponse, *responses.PredictResponse, error) {
	parent, err := t.parentRepo.FindParentByPhoneNumber(req.PhoneNumber)

	if parent == nil {
		return nil, nil, fmt.Errorf("parent dengan nomor telepon %s tidak ditemukan", req.PhoneNumber)
	}

	toddlerMapping := models.Toddler{
		ParentID:   parent.ID,
		Name:       req.Name,
		Birthdate:  req.Birthdate,
		Sex:        req.Sex,
		Height:     req.Height,
		LocationID: req.LocationID,
	}

	if err != nil {
		return nil, nil, err
	}

	toddler, err := t.repo.CreateToddler(&toddlerMapping)

	if err != nil {
		return nil, nil, err
	}

	predict, err := t.predict.CreateIndividualPredict(req, toddler.LocationID, toddler.ID)

	if err != nil {
		return nil, nil, err
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
		return nil, nil, err
	}

	toddlerResponse := responses.ToddlerResponse{
		ID:                toddler.ID,
		ParentID:          toddler.ParentID,
		LocationID:        toddler.LocationID,
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

// CreateToddlerWithParent implements ToddlerService.
func (t *toddlerService) CreateToddlerWithParent(toddlerReq requests.CreateToddlerRequest, parentReq requests.CreateParentRequest) (*responses.ToddlerResponse, *responses.ParentResponse, *responses.PredictResponse, error) {
	parentMapping := models.Parent{
		Name:        parentReq.Name,
		PhoneNumber: parentReq.PhoneNumber,
		Address:     parentReq.Address,
		Nik:         parentReq.Nik,
		Job:         parentReq.Job,
		LocationID:  parentReq.LocationID,
	}

	parent, err := t.parentRepo.CreateParent(&parentMapping)
	if err != nil {
		return nil, nil, nil, err
	}

	toddlerMapping := models.Toddler{
		ParentID:   parent.ID,
		Name:       toddlerReq.Name,
		Birthdate:  toddlerReq.Birthdate,
		Sex:        toddlerReq.Sex,
		Height:     toddlerReq.Height,
		LocationID: toddlerReq.LocationID,
	}

	toddler, err := t.repo.CreateToddler(&toddlerMapping)
	if err != nil {
		return nil, nil, nil, err
	}

	predict, err := t.predict.CreateIndividualPredict(toddlerReq, parentReq.LocationID, toddler.ID)

	if err != nil {
		return nil, nil, nil, err
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

	_, err = t.repo.UpdateToddlerByID(toddler.ID, parentReq.LocationID, &toddlerModel)

	if err != nil {
		return nil, nil, nil, err
	}

	toddler.NutritionalStatus = predict.NutritionalStatus

	toddlerResponse := responses.ToddlerResponse{
		ID:                toddler.ID,
		ParentID:          toddler.ParentID,
		LocationID:        toddler.LocationID,
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
			ID:                v.ID,
			ParentID:          v.ParentID,
			LocationID:        v.Parent.LocationID,
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

	return toddlerResponse, err
}

// GetToddlerByID implements ToddlerService.
func (t *toddlerService) GetToddlerByID(id int, locationID int) (*responses.ToddlerResponse, error) {
	toddler, err := t.repo.GetToddlerByID(id, locationID)

	if err != nil {
		return nil, err
	}

	toddlerResponse := responses.ToddlerResponse{
		ID:                toddler.ID,
		ParentID:          toddler.ParentID,
		LocationID:        toddler.LocationID,
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

// UpdateToddlerByID implements ToddlerService.
func (t *toddlerService) UpdateToddlerByID(
	id int,
	locationID int,
	req requests.UpdateToddlerRequest,
) (*responses.ToddlerResponse, *responses.PredictResponse, error) {

	if req.Name != nil && strings.TrimSpace(*req.Name) == "" {
		return nil, nil, fmt.Errorf("nama tidak boleh kosong")
	}
	if req.Height != nil && *req.Height <= 0 {
		return nil, nil, fmt.Errorf("tinggi badan tidak boleh 0 atau negatif")
	}
	if req.Sex != "" && req.Sex != "M" && req.Sex != "F" {
		return nil, nil, fmt.Errorf("jenis kelamin harus M atau F")
	}
	if req.Birthdate != nil && req.Birthdate.IsZero() {
		return nil, nil, fmt.Errorf("tanggal lahir tidak valid")
	}
	if req.PhoneNumber != nil && strings.TrimSpace(*req.PhoneNumber) == "" {
		return nil, nil, fmt.Errorf("nomor telepon tidak boleh kosong")
	}

	var url string
	var err error

	if req.ProfilePicture != nil && req.ProfilePicture.Filename != "" && req.ProfilePicture.Size > 0 {
		url, err = t.s3.UploadFile(req.ProfilePicture, "toddlers")
		if err != nil {
			return nil, nil, fmt.Errorf("gagal upload foto: %v", err)
		}
	}

	var parentID int
	if req.PhoneNumber != nil {
		parent, err := t.parentRepo.FindParentByPhoneNumber(*req.PhoneNumber)
		if err != nil {
			return nil, nil, fmt.Errorf("orang tua dengan nomor HP %s tidak ditemukan", *req.PhoneNumber)
		}
		parentID = parent.ID
	}

	toddlerMapping := models.Toddler{}
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

	toddler, err := t.repo.UpdateToddlerByID(id, locationID, &toddlerMapping)
	if err != nil {
		return nil, nil, fmt.Errorf("gagal update data toddler: %w", err)
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
		toddlerRequest.LocationID,
		id,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("gagal membuat prediksi: %w", err)
	}

	toddlerResponse := responses.ToddlerResponse{
		ID:             toddler.ID,
		ParentID:       toddler.ParentID,
		LocationID:     toddler.LocationID,
		Name:           toddler.Name,
		Birthdate:      toddler.Birthdate,
		Sex:            toddler.Sex,
		Height:         toddler.Height,
		ProfilePicture: toddler.ProfilePicture,
		CreatedAt:      toddler.CreatedAt,
		UpdatedAt:      toddler.UpdatedAt,
	}

	predictResponse := responses.PredictResponse{
		ID:                predict.ID,
		ToddlerID:         predict.ToddlerID,
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
