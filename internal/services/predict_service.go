package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"grovia/internal/dto/requests"
	"grovia/internal/dto/responses"
	"grovia/internal/models"
	"grovia/internal/repositories"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type PredictService interface {
	CreateIndividualPredict(req requests.CreateToddlerRequest, locationID, toddlerID int) (*responses.PredictResponse, error)
	CreateGroupPredict(filePath string) ([]byte, error)
	GetAllPredict(locationID int) ([]responses.PredictResponse, error)
	GetAllPredictByToddlerID(locationID, toddlerID int) ([]responses.PredictResponse, error)
	GetPredictByID(id int) (*responses.PredictResponse, error)
	UpdatePredictByID(id int, req *requests.UpdatePredictRequest) (*responses.PredictResponse, error)
	DeletePredictByID(id, locationID int) error
	GetAllPredictAllLocation() ([]responses.PredictResponse, error)
}

type predictService struct {
	repo     repositories.PredictRepository
	mlAPIURL string
}

// GetAllPredictAllLocation implements PredictService.
func (p *predictService) GetAllPredictAllLocation() ([]responses.PredictResponse, error) {
	predicts, err := p.repo.GetAllPredictAllLocation()

	if err != nil {
		return nil, err
	}

	var predictResponses []responses.PredictResponse
	for _, v := range predicts {
		predictResponses = append(predictResponses, responses.PredictResponse{
			ID:                v.ID,
			ToddlerID:         v.ToddlerID,
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

	return predictResponses, nil
}

// CreateGroupPredict implements PredictService.
func (p *predictService) CreateGroupPredict(filePath string) ([]byte, error) {
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)

	fileWriter, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return nil, err
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	_, err = io.Copy(fileWriter, file)
	if err != nil {
		return nil, err
	}
	writer.Close()

	req, err := http.NewRequest("POST", p.mlAPIURL+"/predict-group", payload)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("FastAPI error: %s", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// CreateIndividualPredict implements PredictService.
func (p *predictService) CreateIndividualPredict(req requests.CreateToddlerRequest, locationID, toddlerID int) (*responses.PredictResponse, error) {
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
		return nil, err
	}
	defer resp.Body.Close()

	var mlResult map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&mlResult); err != nil {
		return nil, err
	}

	predictModel := &models.Predict{
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
		return nil, err
	}

	return &responses.PredictResponse{
		ID:                saved.ID,
		ToddlerID:         saved.ToddlerID,
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

// DeletePredictByID implements PredictService.
func (p *predictService) DeletePredictByID(id int, locationID int) error {
	return p.repo.DeletePredictByID(id, locationID)
}

// GetAllPredict implements PredictService.
func (p *predictService) GetAllPredict(locationID int) ([]responses.PredictResponse, error) {
	predicts, err := p.repo.GetAllPredict(locationID)
	if err != nil {
		return nil, err
	}

	var responsesList []responses.PredictResponse
	for _, pred := range predicts {
		responsesList = append(responsesList, responses.PredictResponse{
			ID:                pred.ID,
			ToddlerID:         pred.ToddlerID,
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

// GetAllPredictByToddlerID implements PredictService.
func (p *predictService) GetAllPredictByToddlerID(locationID int, toddlerID int) ([]responses.PredictResponse, error) {
	predicts, err := p.repo.GetAllPredictByToddlerID(locationID, toddlerID)
	if err != nil {
		return nil, err
	}

	var responsesList []responses.PredictResponse
	for _, pred := range predicts {
		responsesList = append(responsesList, responses.PredictResponse{
			ID:                pred.ID,
			ToddlerID:         pred.ToddlerID,
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

// GetPredictByID implements PredictService.
func (p *predictService) GetPredictByID(id int) (*responses.PredictResponse, error) {
	predict, err := p.repo.GetPredictByID(id)
	if err != nil {
		return nil, err
	}

	return &responses.PredictResponse{
		ID:                predict.ID,
		ToddlerID:         predict.ToddlerID,
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

// UpdatePredictByID implements PredictService.
func (p *predictService) UpdatePredictByID(id int, req *requests.UpdatePredictRequest) (*responses.PredictResponse, error) {
	predictModel := &models.Predict{
		ID:                id,
		Height:            *req.Height,
		Age:               *req.Age,
		Sex:               req.Sex,
		Zscore:            *req.Zscore,
		NutritionalStatus: *req.NutritionalStatus,
		UpdatedAt:         time.Now(),
	}

	updated, err := p.repo.UpdatePredictByID(id, predictModel)
	if err != nil {
		return nil, err
	}

	return &responses.PredictResponse{
		ID:                updated.ID,
		ToddlerID:         updated.ToddlerID,
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
