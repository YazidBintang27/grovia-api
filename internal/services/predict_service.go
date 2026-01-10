package services

import (
	"bytes"
	"encoding/json"
	"grovia/internal/dto/requests"
	"grovia/internal/dto/responses"
	"grovia/internal/models"
	"grovia/internal/repositories"
	"grovia/pkg"
	"io"
	"math"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type PredictService interface {
	CreateIndividualPredict(req requests.CreateToddlerRequest, locationID, toddlerID, userID int) (*responses.PredictResponse, error)
	CreateGroupPredict(filePath string) ([]byte, error)
	GetAllPredict(locationID int, pageStr, limitStr string) ([]responses.PredictResponse, *responses.PaginationMeta, error)
	GetAllPredictByToddlerID(locationID, toddlerID int) ([]responses.PredictResponse, error)
	GetPredictByID(id int) (*responses.PredictResponse, error)
	UpdatePredictByID(id int, req *requests.UpdatePredictRequest) (*responses.PredictResponse, error)
	DeletePredictByID(id, locationID, userID int) error
	GetAllPredictAllLocation(pageStr, limitStr string) ([]responses.PredictResponse, *responses.PaginationMeta, error)
}

type predictService struct {
	repo     repositories.PredictRepository
	mlAPIURL string
}

func (p *predictService) GetAllPredictAllLocation(pageStr, limitStr string) ([]responses.PredictResponse, *responses.PaginationMeta, error) {
	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 1
	}

	offset := (page - 1) * limit

	predicts, total, err := p.repo.GetAllPredictAllLocation(limit, offset)

	if err != nil {
		return nil, nil, pkg.NewInternalServerError("Gagal mengambil data prediksi")
	}

	totalPage := int(math.Ceil(float64(total) / float64(limit)))

	var predictResponses []responses.PredictResponse
	for _, v := range predicts {
		predictResponses = append(predictResponses, responses.PredictResponse{
			ID:                v.ID,
			ToddlerID:         v.ToddlerID,
			CreatedByID:       v.CreatedByID,
			Name:              v.Name,
			Height:            v.Height,
			Age:               v.Age,
			Sex:               v.Sex,
			Zscore:            v.Zscore,
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

	return predictResponses, &meta, nil
}

func (p *predictService) CreateGroupPredict(filePath string) ([]byte, error) {
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)

	fileWriter, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return nil, pkg.NewInternalServerError("Gagal membuat form file")
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, pkg.NewInternalServerError("Gagal membuka file")
	}
	defer file.Close()

	_, err = io.Copy(fileWriter, file)
	if err != nil {
		return nil, pkg.NewInternalServerError("Gagal menyalin file")
	}
	writer.Close()

	req, err := http.NewRequest("POST", p.mlAPIURL+"/predict-group", payload)
	if err != nil {
		return nil, pkg.NewInternalServerError("Gagal membuat request")
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, pkg.NewInternalServerError("Gagal menghubungi ML API")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, pkg.NewInternalServerError("ML API error: " + resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, pkg.NewInternalServerError("Gagal membaca response")
	}

	return data, nil
}

func (p *predictService) CreateIndividualPredict(req requests.CreateToddlerRequest, locationID, toddlerID, userID int) (*responses.PredictResponse, error) {
	if err := pkg.ValidateStruct(req); err != nil {
		return nil, pkg.NewBadRequestError(err.Error())
	}
	today := time.Now()

	age := (today.Year()-req.Birthdate.Year())*12 + int(today.Month()) - int(req.Birthdate.Month())

	if today.Day() < req.Birthdate.Day() {
		age--
	}
	payload, _ := json.Marshal(map[string]any{
		"height": req.Height,
		"age":    age,
		"gender": req.Sex,
	})

	resp, err := http.Post(p.mlAPIURL+"/predict-individual", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return nil, pkg.NewInternalServerError("Gagal menghubungi ML API")
	}
	defer resp.Body.Close()

	var mlResult map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&mlResult); err != nil {
		return nil, pkg.NewInternalServerError("Gagal decode response ML API")
	}

	predictModel := &models.Predict{
		CreatedByID:       userID,
		DeletedByID:       nil,
		Name:              req.Name,
		Height:            req.Height,
		Age:               age,
		Sex:               req.Sex,
		Zscore:            mlResult["zscore"].(float64),
		NutritionalStatus: mlResult["nutritionalStatus"].(string),
		LocationID:        locationID,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	saved, err := p.repo.CreateIndividualPredict(predictModel, locationID, toddlerID)
	if err != nil {
		return nil, pkg.NewInternalServerError("Gagal menyimpan prediksi")
	}

	return &responses.PredictResponse{
		ID:                saved.ID,
		ToddlerID:         saved.ToddlerID,
		CreatedByID:       saved.CreatedByID,
		Name:              saved.Name,
		Height:            saved.Height,
		Age:               saved.Age,
		Sex:               saved.Sex,
		Zscore:            saved.Zscore,
		NutritionalStatus: saved.NutritionalStatus,
		CreatedAt:         saved.CreatedAt,
		UpdatedAt:         saved.UpdatedAt,
	}, nil
}

func (p *predictService) DeletePredictByID(id int, locationID, userID int) error {
	err := p.repo.DeletePredictByID(id, locationID, userID)
	if err != nil {
		return pkg.NewInternalServerError("Gagal menghapus prediksi")
	}
	return nil
}

func (p *predictService) GetAllPredict(locationID int, pageStr, limitStr string) ([]responses.PredictResponse, *responses.PaginationMeta, error) {
	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 1
	}

	offset := (page - 1) * limit

	predicts, total, err := p.repo.GetAllPredict(locationID, limit, offset)

	if err != nil {
		return nil, nil, pkg.NewInternalServerError("Gagal mengambil data prediksi")
	}

	totalPage := int(math.Ceil(float64(total) / float64(limit)))

	var responsesList []responses.PredictResponse
	for _, pred := range predicts {
		responsesList = append(responsesList, responses.PredictResponse{
			ID:                pred.ID,
			ToddlerID:         pred.ToddlerID,
			CreatedByID:       pred.CreatedByID,
			Name:              pred.Name,
			Height:            pred.Height,
			Age:               pred.Age,
			Sex:               pred.Sex,
			Zscore:            pred.Zscore,
			NutritionalStatus: pred.NutritionalStatus,
			CreatedAt:         pred.CreatedAt,
			UpdatedAt:         pred.UpdatedAt,
		})
	}

	meta := responses.PaginationMeta{
		Page:      page,
		Limit:     limit,
		TotalData: total,
		TotalPage: totalPage,
	}

	return responsesList, &meta, nil
}

func (p *predictService) GetAllPredictByToddlerID(locationID int, toddlerID int) ([]responses.PredictResponse, error) {
	predicts, err := p.repo.GetAllPredictByToddlerID(locationID, toddlerID)
	if err != nil {
		return nil, pkg.NewInternalServerError("Gagal mengambil data prediksi")
	}

	var responsesList []responses.PredictResponse
	for _, pred := range predicts {
		responsesList = append(responsesList, responses.PredictResponse{
			ID:                pred.ID,
			ToddlerID:         pred.ToddlerID,
			CreatedByID:       pred.CreatedByID,
			Name:              pred.Name,
			Height:            pred.Height,
			Age:               pred.Age,
			Sex:               pred.Sex,
			Zscore:            pred.Zscore,
			NutritionalStatus: pred.NutritionalStatus,
			CreatedAt:         pred.CreatedAt,
			UpdatedAt:         pred.UpdatedAt,
		})
	}

	return responsesList, nil
}

func (p *predictService) GetPredictByID(id int) (*responses.PredictResponse, error) {
	predict, err := p.repo.GetPredictByID(id)
	if err != nil {
		return nil, pkg.NewNotFoundError("Prediksi tidak ditemukan")
	}

	return &responses.PredictResponse{
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
	}, nil
}

func (p *predictService) UpdatePredictByID(id int, req *requests.UpdatePredictRequest) (*responses.PredictResponse, error) {
	predictModel := &models.Predict{
		ID:                id,
		DeletedByID:       nil,
		Height:            *req.Height,
		Age:               *req.Age,
		Sex:               *req.Sex,
		Zscore:            *req.Zscore,
		NutritionalStatus: *req.NutritionalStatus,
		UpdatedAt:         time.Now(),
	}

	updated, err := p.repo.UpdatePredictByID(id, predictModel)
	if err != nil {
		return nil, pkg.NewInternalServerError("Gagal update prediksi")
	}

	return &responses.PredictResponse{
		ID:                updated.ID,
		ToddlerID:         updated.ToddlerID,
		CreatedByID:       updated.CreatedByID,
		Name:              updated.Name,
		Height:            updated.Height,
		Age:               updated.Age,
		Sex:               updated.Sex,
		Zscore:            updated.Zscore,
		NutritionalStatus: updated.NutritionalStatus,
		CreatedAt:         updated.CreatedAt,
		UpdatedAt:         updated.UpdatedAt,
	}, nil
}

func NewPredictService(repo repositories.PredictRepository, mlAPIURL string) PredictService {
	return &predictService{repo: repo, mlAPIURL: mlAPIURL}
}