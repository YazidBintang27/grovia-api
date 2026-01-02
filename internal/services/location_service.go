package services

import (
	"context"
	"errors"
	"fmt"
	"grovia/internal/dto/requests"
	"grovia/internal/dto/responses"
	"grovia/internal/models"
	"grovia/internal/repositories"
	"log"
	"math"
	"strconv"
	"strings"
)

type LocationService interface {
	CreateLocation(ctx context.Context, req requests.LocationRequest, userID int) (*responses.LocationResponse, error)
	GetAllLocation(name, pageStr, limitStr string) ([]responses.LocationResponse, *responses.PaginationMeta, error)
	GetLocationByID(id int) (*responses.LocationResponse, error)
	UpdateLocationByID(ctx context.Context, id, userID int, req requests.LocationRequest) (*responses.LocationResponse, error)
	DeleteLocationByID(id, userID int) error
}

type locationService struct {
	repo repositories.LocationRepository
	s3   *S3Service
}

// CreateLocation implements LocationService.
func (l *locationService) CreateLocation(ctx context.Context, req requests.LocationRequest, userID int) (*responses.LocationResponse, error) {
	if strings.TrimSpace(req.Name) == "" {
		return nil, fmt.Errorf("nama lokasi tidak boleh kosong")
	}
	if strings.TrimSpace(req.Address) == "" {
		return nil, fmt.Errorf("alamat lokasi tidak boleh kosong")
	}

	locationMapping := models.Location{
		Name:    req.Name,
		Address: req.Address,
	}

	var url string
	var err error
	if req.Picture != nil {
		url, err = l.s3.UploadFile(ctx, req.Picture, "locations")
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
func (l *locationService) DeleteLocationByID(id, userID int) error {
	return l.repo.DeleteLocationByID(id, userID)
}

// GetAllLocation implements LocationService.
func (l *locationService) GetAllLocation(name, pageStr, limitStr string) ([]responses.LocationResponse, *responses.PaginationMeta, error) {
	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 1
	}

	offset := (page - 1) * limit

	locations, total, err := l.repo.GetAllLocation(name, limit, offset)

	totalPage := int(math.Ceil(float64(total) / float64(limit)))

	if err != nil {
		return nil, nil, err
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

	meta := responses.PaginationMeta{
		Page:      page,
		Limit:     limit,
		TotalData: total,
		TotalPage: totalPage,
	}

	return locationsResponse, &meta, nil
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
func (l *locationService) UpdateLocationByID(ctx context.Context, id, userID int, req requests.LocationRequest) (*responses.LocationResponse, error) {
	if strings.TrimSpace(req.Name) == "" {
		return nil, errors.New("nama lokasi tidak boleh kosong")
	}

	if strings.TrimSpace(req.Address) == "" {
		return nil, errors.New("alamat lokasi tidak boleh kosong")
	}

	var url string
	var err error
	if req.Picture != nil && req.Picture.Filename != "" && req.Picture.Size > 0 {
		url, err = l.s3.UploadFile(ctx, req.Picture, "locations")
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
