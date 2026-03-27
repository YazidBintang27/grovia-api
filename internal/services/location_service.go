package services

import (
	"context"
	"grovia/internal/dto/requests"
	"grovia/internal/dto/responses"
	"grovia/internal/repositories"
	"grovia/pkg"
	"log"
	"math"
	"strconv"
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

func (l *locationService) CreateLocation(ctx context.Context, req requests.LocationRequest, userID int) (*responses.LocationResponse, error) {
	if err := pkg.ValidateStruct(req); err != nil {
		return nil, pkg.NewBadRequestError(err.Error())
	}

	var url string
	var err error
	if req.Picture != nil {
		url, err = l.s3.UploadFile(ctx, req.Picture, "locations")
		if err != nil {
			return nil, pkg.NewInternalServerError("Gagal upload gambar lokasi")
		}
	}

	locationMapping := req.ToModel(url)

	location, err := l.repo.CreateLocation(&locationMapping)

	if err != nil {
		return nil, pkg.NewInternalServerError("Gagal membuat lokasi")
	}

	return responses.FromModelLocation(*location), nil
}

func (l *locationService) DeleteLocationByID(id, userID int) error {
	err := l.repo.DeleteLocationByID(id, userID)
	if err != nil {
		return pkg.NewInternalServerError("Gagal menghapus lokasi")
	}
	return nil
}

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
		return nil, nil, pkg.NewInternalServerError("Gagal mengambil data lokasi")
	}

	meta := responses.PaginationMeta{
		Page:      page,
		Limit:     limit,
		TotalData: total,
		TotalPage: totalPage,
	}

	return responses.FromModelLocationList(locations), &meta, nil
}

func (l *locationService) GetLocationByID(id int) (*responses.LocationResponse, error) {
	location, err := l.repo.GetLocationByID(id)

	if err != nil {
		return nil, pkg.NewNotFoundError("Lokasi tidak ditemukan")
	}

	return responses.FromModelLocation(*location), nil
}

func (l *locationService) UpdateLocationByID(ctx context.Context, id, userID int, req requests.LocationRequest) (*responses.LocationResponse, error) {
	if err := pkg.ValidateStruct(req); err != nil {
		return nil, pkg.NewBadRequestError(err.Error())
	}

	var url string
	var err error
	if req.Picture != nil && req.Picture.Filename != "" && req.Picture.Size > 0 {
		url, err = l.s3.UploadFile(ctx, req.Picture, "locations")
		if err != nil {
			return nil, pkg.NewInternalServerError("Gagal upload gambar lokasi")
		}
	}

	log.Println("[DEBUG] Location Picture URL:", url)

	locationMapping := req.ToModel(url)

	location, err := l.repo.UpdateLocationByID(id, &locationMapping)
	if err != nil {
		return nil, pkg.NewInternalServerError("Gagal update data lokasi")
	}

	return responses.FromModelLocation(*location), nil
}

func NewLocationService(repo repositories.LocationRepository, s3 *S3Service) LocationService {
	return &locationService{repo: repo, s3: s3}
}