package services_test

import (
	"context"
	"errors"
	"grovia/internal/dto/requests"
	"grovia/internal/dto/responses"
	"grovia/internal/mocks"
	"grovia/internal/models"
	"grovia/internal/services"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ─── Helper Setup ────────────────────────────────────────────────────────────

type toddlerServiceDeps struct {
	toddlerRepo *mocks.MockToddlerRepository
	parentRepo  *mocks.MockParentRepository
	predictSvc  *mocks.MockPredictService
	svc         services.ToddlerService
}

func newToddlerServiceDeps() toddlerServiceDeps {
	toddlerRepo := new(mocks.MockToddlerRepository)
	parentRepo := new(mocks.MockParentRepository)
	predictSvc := new(mocks.MockPredictService)

	svc := services.NewToddlerService(toddlerRepo, parentRepo, nil, predictSvc)

	return toddlerServiceDeps{
		toddlerRepo: toddlerRepo,
		parentRepo:  parentRepo,
		predictSvc:  predictSvc,
		svc:         svc,
	}
}

func (d *toddlerServiceDeps) assertAll(t *testing.T) {
	d.toddlerRepo.AssertExpectations(t)
	d.parentRepo.AssertExpectations(t)
	d.predictSvc.AssertExpectations(t)
}

// ─── Fixtures ────────────────────────────────────────────────────────────────

func sampleToddler() *models.Toddler {
	return &models.Toddler{
		ID:                1,
		Name:              "Anak Test",
		LocationID:        1,
		ParentID:          10,
		NutritionalStatus: "Normal",
	}
}

func sampleParent() *models.Parent {
	return &models.Parent{
		ID:          10,
		Name:        "Orang Tua Test",
		PhoneNumber: "08123456789",
		LocationID:  1,
	}
}

func samplePredictResponse() *responses.PredictResponse {
	return &responses.PredictResponse{
		ID:                1,
		NutritionalStatus: "Normal",
		Zscore:            -0.5,
	}
}

func sampleCreateToddlerRequest() requests.CreateToddlerRequest {
	return requests.CreateToddlerRequest{
		Name:              "Anak Test",
		PhoneNumber:       "08123456789",
		Height:            85.0,
		Sex:               "male",
		NutritionalStatus: "normal",
		LocationID:        2,
		Birthdate:         time.Now().AddDate(-2, 0, 0),
	}
}

// ─── GetToddlerByID ───────────────────────────────────────────────────────────

func TestGetToddlerByID_Success(t *testing.T) {
	d := newToddlerServiceDeps()

	d.toddlerRepo.On("GetToddlerByID", 1, 1).Return(sampleToddler(), nil)

	result, err := d.svc.GetToddlerByID(1, 1)

	assert.NoError(t, err)
	assert.Equal(t, "Anak Test", result.Name)
	d.assertAll(t)
}

func TestGetToddlerByID_NotFound(t *testing.T) {
	d := newToddlerServiceDeps()

	d.toddlerRepo.On("GetToddlerByID", 99, 1).Return(nil, errors.New("record not found"))

	result, err := d.svc.GetToddlerByID(99, 1)

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "tidak ditemukan")
	d.assertAll(t)
}

// ─── GetAllToddler ────────────────────────────────────────────────────────────

func TestGetAllToddler_Success(t *testing.T) {
	d := newToddlerServiceDeps()

	toddlers := []models.Toddler{*sampleToddler(), {ID: 2, Name: "Anak Dua"}}
	d.toddlerRepo.On("GetAllToddler", 1, 10, 0, "").Return(toddlers, 2, nil)

	results, meta, err := d.svc.GetAllToddler(1, "", "1", "10")

	assert.NoError(t, err)
	assert.Len(t, results, 2)
	assert.Equal(t, 2, meta.TotalData)
	assert.Equal(t, 1, meta.TotalPage)
	d.assertAll(t)
}

func TestGetAllToddler_WithNameFilter(t *testing.T) {
	d := newToddlerServiceDeps()

	toddlers := []models.Toddler{*sampleToddler()}
	d.toddlerRepo.On("GetAllToddler", 1, 10, 0, "anak").Return(toddlers, 1, nil)

	results, meta, err := d.svc.GetAllToddler(1, "anak", "1", "10")

	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, 1, meta.TotalData)
	d.assertAll(t)
}

func TestGetAllToddler_RepoError(t *testing.T) {
	d := newToddlerServiceDeps()

	d.toddlerRepo.On("GetAllToddler", 1, 10, 0, "").Return([]models.Toddler{}, 0, errors.New("db error"))

	results, meta, err := d.svc.GetAllToddler(1, "", "1", "10")

	assert.Nil(t, results)
	assert.Nil(t, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Gagal mengambil")
	d.assertAll(t)
}

func TestGetAllToddler_DefaultPagination(t *testing.T) {
	d := newToddlerServiceDeps()

	// page=0, limit=0 → default ke page=1, limit=1
	d.toddlerRepo.On("GetAllToddler", 1, 1, 0, "").Return([]models.Toddler{}, 0, nil)

	_, meta, err := d.svc.GetAllToddler(1, "", "0", "0")

	assert.NoError(t, err)
	assert.Equal(t, 1, meta.Page)
	assert.Equal(t, 1, meta.Limit)
	d.assertAll(t)
}

// ─── GetAllToddlerAllLocation ─────────────────────────────────────────────────

func TestGetAllToddlerAllLocation_Success(t *testing.T) {
	d := newToddlerServiceDeps()

	toddlers := []models.Toddler{*sampleToddler()}
	d.toddlerRepo.On("GetAllToddlerAllLocation", "", 10, 0).Return(toddlers, 1, nil)

	results, meta, err := d.svc.GetAllToddlerAllLocation("", "1", "10")

	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, 1, meta.TotalData)
	d.assertAll(t)
}

func TestGetAllToddlerAllLocation_RepoError(t *testing.T) {
	d := newToddlerServiceDeps()

	d.toddlerRepo.On("GetAllToddlerAllLocation", "", 10, 0).Return([]models.Toddler{}, 0, errors.New("db error"))

	results, meta, err := d.svc.GetAllToddlerAllLocation("", "1", "10")

	assert.Nil(t, results)
	assert.Nil(t, meta)
	assert.Error(t, err)
	d.assertAll(t)
}

// ─── DeleteToddlerByID ────────────────────────────────────────────────────────

func TestDeleteToddlerByID_Success(t *testing.T) {
	d := newToddlerServiceDeps()

	d.toddlerRepo.On("DeleteToddlerByID", 1, 1, 10).Return(nil)

	err := d.svc.DeleteToddlerByID(1, 1, 10)

	assert.NoError(t, err)
	d.assertAll(t)
}

func TestDeleteToddlerByID_Failed(t *testing.T) {
	d := newToddlerServiceDeps()

	d.toddlerRepo.On("DeleteToddlerByID", 99, 1, 10).Return(errors.New("not found"))

	err := d.svc.DeleteToddlerByID(99, 1, 10)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Gagal menghapus")
	d.assertAll(t)
}

// ─── CheckToddlerExists ───────────────────────────────────────────────────────

func TestCheckToddlerExists_Found(t *testing.T) {
	d := newToddlerServiceDeps()

	d.parentRepo.On("FindParentByPhoneNumber", "08123456789").Return(sampleParent(), nil)
	d.toddlerRepo.On("FindToddlerByName", 10, "Anak Test").Return(true, sampleToddler(), nil)

	exists, toddler, err := d.svc.CheckToddlerExists("08123456789", "Anak Test")

	assert.NoError(t, err)
	assert.True(t, exists)
	assert.NotNil(t, toddler)
	d.assertAll(t)
}

func TestCheckToddlerExists_ParentNotFound(t *testing.T) {
	d := newToddlerServiceDeps()

	d.parentRepo.On("FindParentByPhoneNumber", "0000").Return(nil, errors.New("not found"))

	exists, toddler, err := d.svc.CheckToddlerExists("0000", "Anak Test")

	assert.False(t, exists)
	assert.Nil(t, toddler)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Parent tidak ditemukan")
	d.assertAll(t)
}

func TestCheckToddlerExists_ToddlerNotFound(t *testing.T) {
	d := newToddlerServiceDeps()

	d.parentRepo.On("FindParentByPhoneNumber", "08123456789").Return(sampleParent(), nil)
	d.toddlerRepo.On("FindToddlerByName", 10, "Tidak Ada").Return(false, nil, errors.New("not found"))

	exists, toddler, err := d.svc.CheckToddlerExists("08123456789", "Tidak Ada")

	assert.False(t, exists)
	assert.Nil(t, toddler)
	assert.Error(t, err)
	d.assertAll(t)
}

// ─── CreateToddler ────────────────────────────────────────────────────────────

func TestCreateToddler_Success(t *testing.T) {
	d := newToddlerServiceDeps()

	req := sampleCreateToddlerRequest()
	parent := sampleParent()
	toddler := sampleToddler()
	predictResp := samplePredictResponse()

	d.parentRepo.On("FindParentByPhoneNumber", "08123456789").Return(parent, nil)
	d.toddlerRepo.On("CreateToddler", mock.AnythingOfType("*models.Toddler")).Return(toddler, nil)
	d.predictSvc.On("CreateIndividualPredict", req, toddler.LocationID, toddler.ID, 1).Return(predictResp, nil)
	d.toddlerRepo.On("UpdateToddlerByID", toddler.ID, parent.LocationID, mock.AnythingOfType("*models.Toddler")).Return(toddler, nil)

	toddlerResp, predictResult, err := d.svc.CreateToddler(req, 1)

	assert.NoError(t, err)
	assert.NotNil(t, toddlerResp)
	assert.NotNil(t, predictResult)
	assert.Equal(t, "Normal", predictResult.NutritionalStatus)
	d.assertAll(t)
}

func TestCreateToddler_ParentNotFound(t *testing.T) {
	d := newToddlerServiceDeps()

	req := sampleCreateToddlerRequest()
	d.parentRepo.On("FindParentByPhoneNumber", "08123456789").Return(nil, errors.New("not found"))

	toddlerResp, predictResult, err := d.svc.CreateToddler(req, 1)

	assert.Nil(t, toddlerResp)
	assert.Nil(t, predictResult)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "tidak ditemukan")
	d.assertAll(t)
}

func TestCreateToddler_CreateToddlerFailed(t *testing.T) {
	d := newToddlerServiceDeps()

	req := sampleCreateToddlerRequest()
	parent := sampleParent()

	d.parentRepo.On("FindParentByPhoneNumber", "08123456789").Return(parent, nil)
	d.toddlerRepo.On("CreateToddler", mock.AnythingOfType("*models.Toddler")).Return(nil, errors.New("db error"))

	toddlerResp, predictResult, err := d.svc.CreateToddler(req, 1)

	assert.Nil(t, toddlerResp)
	assert.Nil(t, predictResult)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Gagal membuat toddler")
	d.assertAll(t)
}

func TestCreateToddler_PredictFailed(t *testing.T) {
	d := newToddlerServiceDeps()

	req := sampleCreateToddlerRequest()
	parent := sampleParent()
	toddler := sampleToddler()

	d.parentRepo.On("FindParentByPhoneNumber", "08123456789").Return(parent, nil)
	d.toddlerRepo.On("CreateToddler", mock.AnythingOfType("*models.Toddler")).Return(toddler, nil)
	d.predictSvc.On("CreateIndividualPredict", req, toddler.LocationID, toddler.ID, 1).Return(nil, errors.New("ml error"))

	toddlerResp, predictResult, err := d.svc.CreateToddler(req, 1)

	assert.Nil(t, toddlerResp)
	assert.Nil(t, predictResult)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Gagal membuat prediksi")
	d.assertAll(t)
}

func TestCreateToddler_UpdateStatusFailed(t *testing.T) {
	d := newToddlerServiceDeps()

	req := sampleCreateToddlerRequest()
	parent := sampleParent()
	toddler := sampleToddler()
	predictResp := samplePredictResponse()

	d.parentRepo.On("FindParentByPhoneNumber", "08123456789").Return(parent, nil)
	d.toddlerRepo.On("CreateToddler", mock.AnythingOfType("*models.Toddler")).Return(toddler, nil)
	d.predictSvc.On("CreateIndividualPredict", req, toddler.LocationID, toddler.ID, 1).Return(predictResp, nil)
	d.toddlerRepo.On("UpdateToddlerByID", toddler.ID, parent.LocationID, mock.AnythingOfType("*models.Toddler")).Return(nil, errors.New("db error"))

	toddlerResp, predictResult, err := d.svc.CreateToddler(req, 1)

	assert.Nil(t, toddlerResp)
	assert.Nil(t, predictResult)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Gagal update nutritional status")
	d.assertAll(t)
}

// ─── CreateToddlerWithParent ──────────────────────────────────────────────────

func TestCreateToddlerWithParent_Success(t *testing.T) {
	d := newToddlerServiceDeps()

	toddlerReq := sampleCreateToddlerRequest()
	parentReq := requests.CreateParentRequest{
		Name:        "Orang Tua Test",
		PhoneNumber: "08123456789",
		LocationID:  1,
	}

	parent := sampleParent()
	toddler := sampleToddler()
	predictResp := samplePredictResponse()

	d.parentRepo.On("CreateParent", mock.AnythingOfType("*models.Parent")).Return(parent, nil)
	d.toddlerRepo.On("CreateToddler", mock.AnythingOfType("*models.Toddler")).Return(toddler, nil)
	d.predictSvc.On("CreateIndividualPredict", toddlerReq, parentReq.LocationID, toddler.ID, 1).Return(predictResp, nil)
	d.toddlerRepo.On("UpdateToddlerByID", toddler.ID, parentReq.LocationID, mock.AnythingOfType("*models.Toddler")).Return(toddler, nil)

	toddlerResp, parentResp, predictResult, err := d.svc.CreateToddlerWithParent(toddlerReq, parentReq, 1)

	assert.NoError(t, err)
	assert.NotNil(t, toddlerResp)
	assert.NotNil(t, parentResp)
	assert.NotNil(t, predictResult)
	assert.Equal(t, "Normal", predictResult.NutritionalStatus)
	d.assertAll(t)
}

func TestCreateToddlerWithParent_CreateParentFailed(t *testing.T) {
	d := newToddlerServiceDeps()

	toddlerReq := sampleCreateToddlerRequest()
	parentReq := requests.CreateParentRequest{
		Name:        "Orang Tua Test",
		PhoneNumber: "08123456789",
		LocationID:  1,
	}

	d.parentRepo.On("CreateParent", mock.AnythingOfType("*models.Parent")).Return(nil, errors.New("db error"))

	toddlerResp, parentResp, predictResult, err := d.svc.CreateToddlerWithParent(toddlerReq, parentReq, 1)

	assert.Nil(t, toddlerResp)
	assert.Nil(t, parentResp)
	assert.Nil(t, predictResult)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Gagal membuat parent")
	d.assertAll(t)
}

func TestCreateToddlerWithParent_CreateToddlerFailed(t *testing.T) {
	d := newToddlerServiceDeps()

	toddlerReq := sampleCreateToddlerRequest()
	parentReq := requests.CreateParentRequest{
		Name:        "Orang Tua Test",
		PhoneNumber: "08123456789",
		LocationID:  1,
	}

	d.parentRepo.On("CreateParent", mock.AnythingOfType("*models.Parent")).Return(sampleParent(), nil)
	d.toddlerRepo.On("CreateToddler", mock.AnythingOfType("*models.Toddler")).Return(nil, errors.New("db error"))

	toddlerResp, parentResp, predictResult, err := d.svc.CreateToddlerWithParent(toddlerReq, parentReq, 1)

	assert.Nil(t, toddlerResp)
	assert.Nil(t, parentResp)
	assert.Nil(t, predictResult)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Gagal membuat toddler")
	d.assertAll(t)
}

// ─── UpdateToddlerByID ────────────────────────────────────────────────────────

func TestUpdateToddlerByID_Success(t *testing.T) {
	d := newToddlerServiceDeps()

	phone := "08123456789"
	req := requests.UpdateToddlerRequest{
		Name:        strPtr("Anak Update"),
		Height:      float64Ptr(90.0),
		PhoneNumber: &phone,
	}

	predictResp := samplePredictResponse()
	toddler := sampleToddler()

	// UpdateToddlerByID memanggil predict dulu, lalu lookup parent, lalu update
	d.predictSvc.On("CreateIndividualPredict", req.ToCreateRequest(), 1, 1, 10).Return(predictResp, nil)
	d.parentRepo.On("FindParentByPhoneNumber", phone).Return(sampleParent(), nil)
	d.toddlerRepo.On("UpdateToddlerByID", 1, 1, mock.AnythingOfType("*models.Toddler")).Return(toddler, nil)

	toddlerResp, predictResult, err := d.svc.UpdateToddlerByID(context.Background(), 1, 1, 10, req)

	assert.NoError(t, err)
	assert.NotNil(t, toddlerResp)
	assert.NotNil(t, predictResult)
	d.assertAll(t)
}

func TestUpdateToddlerByID_PredictFailed(t *testing.T) {
	d := newToddlerServiceDeps()

	phone := "08123456789"
	req := requests.UpdateToddlerRequest{
		Name:        strPtr("Anak Update"),
		PhoneNumber: &phone,
	}

	d.predictSvc.On("CreateIndividualPredict", req.ToCreateRequest(), 1, 1, 10).Return(nil, errors.New("ml error"))

	toddlerResp, predictResult, err := d.svc.UpdateToddlerByID(context.Background(), 1, 1, 10, req)

	assert.Nil(t, toddlerResp)
	assert.Nil(t, predictResult)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Gagal membuat prediksi")
	d.assertAll(t)
}

func TestUpdateToddlerByID_ParentNotFound(t *testing.T) {
	d := newToddlerServiceDeps()

	phone := "0000"
	req := requests.UpdateToddlerRequest{
		Name:        strPtr("Anak Update"),
		PhoneNumber: &phone,
	}

	predictResp := samplePredictResponse()

	d.predictSvc.On("CreateIndividualPredict", req.ToCreateRequest(), 1, 1, 10).Return(predictResp, nil)
	d.parentRepo.On("FindParentByPhoneNumber", phone).Return(nil, errors.New("not found"))

	toddlerResp, predictResult, err := d.svc.UpdateToddlerByID(context.Background(), 1, 1, 10, req)

	assert.Nil(t, toddlerResp)
	assert.Nil(t, predictResult)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "tidak ditemukan")
	d.assertAll(t)
}

// ─── UpdateToddlerByIDWithoutPredict ─────────────────────────────────────────

func TestUpdateToddlerByIDWithoutPredict_Success(t *testing.T) {
	d := newToddlerServiceDeps()

	phone := "08123456789"
	req := requests.UpdateToddlerRequest{
		Name:        strPtr("Anak Update"),
		PhoneNumber: &phone,
	}

	d.parentRepo.On("FindParentByPhoneNumber", phone).Return(sampleParent(), nil)
	d.toddlerRepo.On("UpdateToddlerByID", 1, 1, mock.AnythingOfType("*models.Toddler")).Return(sampleToddler(), nil)

	result, err := d.svc.UpdateToddlerByIDWithoutPredict(context.Background(), 1, 1, 10, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	d.assertAll(t)
}

func TestUpdateToddlerByIDWithoutPredict_ParentNotFound(t *testing.T) {
	d := newToddlerServiceDeps()

	phone := "0000"
	req := requests.UpdateToddlerRequest{
		Name:        strPtr("Anak Update"),
		PhoneNumber: &phone,
	}

	d.parentRepo.On("FindParentByPhoneNumber", phone).Return(nil, errors.New("not found"))

	result, err := d.svc.UpdateToddlerByIDWithoutPredict(context.Background(), 1, 1, 10, req)

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "tidak ditemukan")
	d.assertAll(t)
}

func TestUpdateToddlerByIDWithoutPredict_UpdateFailed(t *testing.T) {
	d := newToddlerServiceDeps()

	phone := "08123456789"
	req := requests.UpdateToddlerRequest{
		Name:        strPtr("Anak Update"),
		PhoneNumber: &phone,
	}

	d.parentRepo.On("FindParentByPhoneNumber", phone).Return(sampleParent(), nil)
	d.toddlerRepo.On("UpdateToddlerByID", 1, 1, mock.AnythingOfType("*models.Toddler")).Return(nil, errors.New("db error"))

	result, err := d.svc.UpdateToddlerByIDWithoutPredict(context.Background(), 1, 1, 10, req)

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Gagal update")
	d.assertAll(t)
}

// ─── Test Helpers ─────────────────────────────────────────────────────────────

func strPtr(s string) *string       { return &s }
func float64Ptr(f float64) *float64 { return &f }
