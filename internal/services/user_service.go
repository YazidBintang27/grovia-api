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

type UserService interface {
	CreateUser(ctx context.Context, req requests.CreateUserRequest, createdBy string, locationID int) (*responses.UserResponse, error)
	GetCurrentUser(id int) (*responses.UserResponse, error)
	GetUsersByRole(requesterRole, name, pageStr, limitStr string, locationID int) ([]responses.UserResponse, *responses.PaginationMeta, error)
	GetUserById(targetUserID int, accesorRole string) (*responses.UserResponse, error)
	UpdateCurrentUser(ctx context.Context, id int, req requests.UpdateUserRequest) (*responses.UserResponse, error)
	UpdateUserByID(ctx context.Context, targetUserID int, req requests.UpdateUserRequest, updaterRole string) (*responses.UserResponse, error)
	DeleteCurrentUser(id int) error
	DeleteUserByID(targetUserID int, role string) error
}

type userService struct {
	repo repositories.UserRepository
	s3   *S3Service
}

// GetUserById implements UserService.
func (u *userService) GetUserById(targetUserID int, accesorRole string) (*responses.UserResponse, error) {
	targetRole, err := u.repo.FindRoleById(targetUserID)

	if err != nil {
		return nil, pkg.NewNotFoundError("User tidak ditemukan")
	}

	if targetRole == "" {
		return nil, pkg.NewNotFoundError("User tidak ditemukan")
	}

	if !u.rolePermission(accesorRole, targetRole) {
		return nil, pkg.NewForbiddenError("Tidak memiliki akses untuk melihat user dengan role " + targetRole)
	}

	user, err := u.repo.GetUser(targetUserID)

	if err != nil {
		return nil, pkg.NewNotFoundError("User tidak ditemukan")
	}

	response := &responses.UserResponse{
		ID:             user.ID,
		LocationID:     user.LocationID,
		Name:           user.Name,
		PhoneNumber:    user.PhoneNumber,
		Address:        user.Address,
		Nik:            user.Nik,
		ProfilePicture: user.ProfilePicture,
		Role:           user.Role,
		IsActive:       user.IsActive,
		CreatedBy:      user.CreatedBy,
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
	}

	return response, nil
}

// GetUsersByRole implements UserService.
func (u *userService) GetUsersByRole(requesterRole, name, pageStr, limitStr string, locationID int) ([]responses.UserResponse, *responses.PaginationMeta, error) {
	var roles []string
	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	switch requesterRole {
	case pkg.RoleAdmin:
		roles = []string{pkg.RoleKepalaPosyandu, pkg.RoleKader}
	case pkg.RoleKepalaPosyandu:
		roles = []string{pkg.RoleKader}
	default:
		return nil, nil, pkg.NewForbiddenError("Role tidak diizinkan")
	}

	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 1
	}

	offset := (page - 1) * limit

	users, total, err := u.repo.FindUsersByRole(roles, name, locationID, limit, offset)
	if err != nil {
		return nil, nil, pkg.NewNotFoundError("User tidak ditemukan")
	}

	totalPage := int(math.Ceil(float64(total) / float64(limit)))

	var result []responses.UserResponse
	for _, v := range users {
		result = append(result, responses.UserResponse{
			ID:             v.ID,
			LocationID:     v.LocationID,
			Name:           v.Name,
			PhoneNumber:    v.PhoneNumber,
			Address:        v.Address,
			Nik:            v.Nik,
			ProfilePicture: v.ProfilePicture,
			Role:           v.Role,
			IsActive:       v.IsActive,
			CreatedBy:      v.CreatedBy,
			CreatedAt:      v.CreatedAt,
			UpdatedAt:      v.UpdatedAt,
		})
	}

	meta := responses.PaginationMeta{
		Page:      page,
		Limit:     limit,
		TotalData: total,
		TotalPage: totalPage,
	}

	return result, &meta, nil
}

// DeleteUserByID implements UserService.
func (u *userService) DeleteUserByID(targetUserID int, updaterRole string) error {
	targetRole, err := u.repo.FindRoleById(targetUserID)

	if err != nil {
		return pkg.NewNotFoundError("User tidak ditemukan")
	}

	if !u.rolePermission(updaterRole, targetRole) {
		return pkg.NewForbiddenError("Tidak memiliki akses menghapus user dengan Role %s" + targetRole)
	}

	return u.repo.DeleteUser(targetUserID)
}

// CreateUser implements UserService.
func (u *userService) CreateUser(ctx context.Context, req requests.CreateUserRequest, createdBy string, locationID int) (*responses.UserResponse, error) {

	if !u.rolePermission(createdBy, req.Role) {
		return nil, pkg.NewForbiddenError("Tidak memiliki akses membuat user dengan Role %s" + req.Role)
	}

	if err := pkg.ValidateStruct(req); err != nil {
		return nil, pkg.NewBadRequestError(err.Error())
	}

	var url string
	var err error

	if req.ProfilePicture != nil && req.ProfilePicture.Filename != "" && req.ProfilePicture.Size > 0 {
		url, err = u.s3.UploadFile(ctx, req.ProfilePicture, "users")
		if err != nil {
			return nil, pkg.NewInternalServerError("Gagal upload foto profil")
		}
	}

	hashedPassword, err := pkg.HashPassword(req.Password)

	if err != nil {
		return nil, pkg.NewInternalServerError("Gagal memproses password")
	}

	var location int

	if createdBy == "admin" {
		if req.LocationID == 0 {
			return nil, pkg.NewBadRequestError("Admin harus menyertakan location ID")
		}
		location = req.LocationID
	} else {
		if locationID == 0 {
			return nil, pkg.NewBadRequestError("Location ID tidak ditemukan")
		}
		location = locationID
	}

	userMapping := models.User{
		LocationID:     location,
		Name:           req.Name,
		PhoneNumber:    req.PhoneNumber,
		Address:        req.Address,
		Nik:            req.Nik,
		ProfilePicture: url,
		Role:           req.Role,
		IsActive:       true,
		Password:       hashedPassword,
		CreatedBy:      createdBy,
	}

	user, err := u.repo.CreateUser(&userMapping)

	if err != nil {
		return nil, pkg.NewInternalServerError("Gagal membuat user")
	}

	response := responses.UserResponse{
		ID:             user.ID,
		LocationID:     location,
		Name:           user.Name,
		PhoneNumber:    user.PhoneNumber,
		Address:        user.Address,
		Nik:            user.Nik,
		ProfilePicture: user.ProfilePicture,
		Role:           user.Role,
		IsActive:       user.IsActive,
		CreatedBy:      user.CreatedBy,
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
	}

	return &response, nil
}

// DeleteUser implements UserService.
func (u *userService) DeleteCurrentUser(id int) error {
	if err := u.repo.DeleteUser(id); err != nil {
		return pkg.NewInternalServerError("Gagal menghapus user")
	}
	return nil
}

// GetUser implements UserService.
func (u *userService) GetCurrentUser(id int) (*responses.UserResponse, error) {
	user, err := u.repo.GetUser(id)

	if err != nil {
		return nil, pkg.NewNotFoundError("User tidak ditemukan")
	}

	response := &responses.UserResponse{
		ID:             user.ID,
		LocationID:     user.LocationID,
		Name:           user.Name,
		PhoneNumber:    user.PhoneNumber,
		Address:        user.Address,
		Nik:            user.Nik,
		ProfilePicture: user.ProfilePicture,
		Role:           user.Role,
		IsActive:       user.IsActive,
		CreatedBy:      user.CreatedBy,
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
	}

	return response, nil
}

// UpdateCurrentUser implements UserService.
func (u *userService) UpdateCurrentUser(ctx context.Context, id int, req requests.UpdateUserRequest) (*responses.UserResponse, error) {
	if err := pkg.ValidateStruct(req); err != nil {
		return nil, pkg.NewBadRequestError(err.Error())
	}

	var url string
	var err error

	if req.ProfilePicture != nil && req.ProfilePicture.Filename != "" && req.ProfilePicture.Size > 0 {
		url, err = u.s3.UploadFile(ctx, req.ProfilePicture, "users")
		if err != nil {
			return nil, pkg.NewInternalServerError("Gagal upload foto profil")
		}
	}

	userMapping := models.User{}
	if req.Name != nil {
		userMapping.Name = *req.Name
	}
	if req.PhoneNumber != nil {
		userMapping.PhoneNumber = *req.PhoneNumber
	}
	if req.Address != nil {
		userMapping.Address = *req.Address
	}
	if req.Nik != nil {
		userMapping.Nik = *req.Nik
	}
	if url != "" {
		userMapping.ProfilePicture = url
	}
	if req.Password != nil {
		hashedPassword, err := pkg.HashPassword(*req.Password)
		if err != nil {
			return nil, pkg.NewInternalServerError("Gagal memproses password")
		}
		userMapping.Password = hashedPassword
	}

	user, err := u.repo.UpdateUser(id, &userMapping)
	if err != nil {
		return nil, pkg.NewInternalServerError("Gagal update user")
	}

	return &responses.UserResponse{
		ID:             user.ID,
		LocationID:     user.LocationID,
		Name:           user.Name,
		PhoneNumber:    user.PhoneNumber,
		Address:        user.Address,
		Nik:            user.Nik,
		ProfilePicture: user.ProfilePicture,
		Role:           user.Role,
		CreatedBy:      user.CreatedBy,
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
	}, nil
}

// UpdateUserByID implements UserService.
func (u *userService) UpdateUserByID(
	ctx context.Context,
	targetUserID int,
	req requests.UpdateUserRequest,
	updaterRole string,
) (*responses.UserResponse, error) {

	if err := pkg.ValidateStruct(req); err != nil {
		return nil, pkg.NewBadRequestError(err.Error())
	}

	if req.Role != nil {
		role, err := u.repo.FindRoleById(targetUserID)
		if err != nil {
			return nil, err
		}
		if !u.rolePermission(updaterRole, role) {
			return nil, pkg.NewForbiddenError("Tidak memiliki akses update user dengan Role %s" + role)
		}
	}

	var url string
	var err error

	if req.ProfilePicture != nil && req.ProfilePicture.Filename != "" && req.ProfilePicture.Size > 0 {
		url, err = u.s3.UploadFile(ctx, req.ProfilePicture, "users")
		if err != nil {
			return nil, pkg.NewInternalServerError("Gagal upload foto profil")
		}
	}

	userMapping := models.User{}
	if req.Name != nil {
		userMapping.Name = *req.Name
	}
	if req.LocationID != nil {
		userMapping.LocationID = *req.LocationID
	}
	if req.PhoneNumber != nil {
		userMapping.PhoneNumber = *req.PhoneNumber
	}
	if req.Address != nil {
		userMapping.Address = *req.Address
	}
	if req.Nik != nil {
		userMapping.Nik = *req.Nik
	}
	if req.Role != nil {
		userMapping.Role = *req.Role
	}
	if url != "" {
		userMapping.ProfilePicture = url
	}
	if req.Password != nil {
		hashedPassword, err := pkg.HashPassword(*req.Password)
		if err != nil {
			return nil, pkg.NewInternalServerError("Gagal memproses password")
		}
		userMapping.Password = hashedPassword
	}

	user, err := u.repo.UpdateUser(targetUserID, &userMapping)
	if err != nil {
		return nil, pkg.NewInternalServerError("Gagal update user")
	}

	return &responses.UserResponse{
		ID:             user.ID,
		LocationID:     user.LocationID,
		Name:           user.Name,
		PhoneNumber:    user.PhoneNumber,
		Address:        user.Address,
		Nik:            user.Nik,
		ProfilePicture: user.ProfilePicture,
		Role:           user.Role,
		CreatedBy:      user.CreatedBy,
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
	}, nil
}

func (u *userService) rolePermission(creatorRole, targetRole string) bool {
	switch creatorRole {
	case pkg.RoleAdmin:
		return true
	case pkg.RoleKepalaPosyandu:
		return targetRole == pkg.RoleKader
	default:
		return false
	}
}

func NewUserService(repo repositories.UserRepository, s3 *S3Service) UserService {
	return &userService{repo: repo, s3: s3}
}
