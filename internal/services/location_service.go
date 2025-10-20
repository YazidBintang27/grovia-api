package services

import (
	"errors"
	"fmt"
	"grovia/internal/dto/requests"
	"grovia/internal/dto/responses"
	"grovia/internal/models"
	"grovia/internal/repositories"
	"log"
	"strings"
)

type LocationService interface {
	CreateLocation(req requests.LocationRequest) (*responses.LocationResponse, error)
	GetAllLocation() ([]responses.LocationResponse, error)
	GetLocationByID(id int) (*responses.LocationResponse, error)
	UpdateLocationByID(id int, req requests.LocationRequest) (*responses.LocationResponse, error)
	DeleteLocationByID(id int) error
}

type locationService struct {
	repo repositories.LocationRepository
	s3   *S3Service
}

// CreateLocation implements LocationService.
func (l *locationService) CreateLocation(req requests.LocationRequest) (*responses.LocationResponse, error) {
	locationMapping := models.Location{
		Name:    req.Name,
		Address: req.Address,
	}

	var url string
	var err error
	if req.Picture != nil {
		url, err = l.s3.UploadFile(req.Picture, "locations")
		if err != nil {
			return nil, err
		}
	}

	if url != "" {
		locationMapping.Picture = url
	}

	location, err := l.repo.CreateLocation(&locationMapping)

	if err != nil {
		return nil, err
	}

	locationResponse := responses.LocationResponse{
		ID:        location.ID,
		Name:      location.Name,
		Address:   location.Address,
		Picture:   location.Picture,
		CreatedAt: location.CreatedAt,
		UpdatedAt: location.UpdatedAt,
	}

	return &locationResponse, nil
}

// DeleteLocationByID implements LocationService.
func (l *locationService) DeleteLocationByID(id int) error {
	return l.repo.DeleteLocationByID(id)
}

// GetAllLocation implements LocationService.
func (l *locationService) GetAllLocation() ([]responses.LocationResponse, error) {
	locations, err := l.repo.GetAllLocation()

	if err != nil {
		return nil, err
	}

	var locationsResponse []responses.LocationResponse

	for _, v := range locations {
		locationsResponse = append(locationsResponse, responses.LocationResponse{
			ID:        v.ID,
			Name:      v.Name,
			Address:   v.Address,
			Picture:   v.Picture,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
		})
	}

	return locationsResponse, nil
}

// GetLocationByID implements LocationService.
func (l *locationService) GetLocationByID(id int) (*responses.LocationResponse, error) {
	location, err := l.repo.GetLocationByID(id)

	if err != nil {
		return nil, err
	}

	locationResponse := responses.LocationResponse{
		ID:        location.ID,
		Name:      location.Name,
		Address:   location.Address,
		Picture:   location.Picture,
		CreatedAt: location.CreatedAt,
		UpdatedAt: location.UpdatedAt,
	}

	return &locationResponse, nil
}

// UpdateLocationByID implements LocationService.
func (l *locationService) UpdateLocationByID(id int, req requests.LocationRequest) (*responses.LocationResponse, error) {
	if strings.TrimSpace(req.Name) == "" {
		return nil, errors.New("nama lokasi tidak boleh kosong")
	}

	if strings.TrimSpace(req.Address) == "" {
		return nil, errors.New("alamat lokasi tidak boleh kosong")
	}

	var url string
	var err error
	if req.Picture != nil && req.Picture.Filename != "" && req.Picture.Size > 0 {
		url, err = l.s3.UploadFile(req.Picture, "locations")
		if err != nil {
			return nil, fmt.Errorf("gagal mengunggah gambar lokasi: %w", err)
		}
	}

	log.Println("[DEBUG] Location Picture URL:", url)

	locationMapping := models.Location{
		Name:    req.Name,
		Address: req.Address,
	}

	if url != "" {
		locationMapping.Picture = url
	}

	location, err := l.repo.UpdateLocationByID(id, &locationMapping)
	if err != nil {
		return nil, fmt.Errorf("gagal memperbarui data lokasi: %w", err)
	}

	locationResponse := responses.LocationResponse{
		ID:        location.ID,
		Name:      location.Name,
		Address:   location.Address,
		Picture:   location.Picture,
		CreatedAt: location.CreatedAt,
		UpdatedAt: location.UpdatedAt,
	}

	return &locationResponse, nil
}

func NewLocationService(repo repositories.LocationRepository, s3 *S3Service) LocationService {
	return &locationService{repo: repo, s3: s3}
}
