package services

import (
	"errors"
	"fmt"
	"grovia/internal/dto/requests"
	"grovia/internal/dto/responses"
	"grovia/internal/models"
	"grovia/internal/repositories"
	"grovia/pkg"
)

type UserService interface {
	CreateUser(req requests.CreateUserRequest, createdBy string, locationID *int) (*responses.UserResponse, error)
	GetCurrentUser(id int) (*responses.UserResponse, error)
	GetUsersByRole(requesterRole string) ([]responses.UserResponse, error)
	GetUserById(targetUserID int, accesorRole string) (*responses.UserResponse, error)
	UpdateCurrentUser(id int, req requests.UpdateUserRequest) (*responses.UserResponse, error)
	UpdateUserByID(targetUserID int, req requests.UpdateUserRequest, updaterRole string) (*responses.UserResponse, error)
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
		return nil, err
	}

	if !u.rolePermission(targetRole, accesorRole) {
		return nil, fmt.Errorf("%s tidak boleh mengubah role ke %s", accesorRole, targetRole)
	}

	user, err := u.repo.GetUser(targetUserID)

	if err != nil {
		return nil, err
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
		CreatedBy:      user.CreatedBy,
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
	}

	return response, nil
}

// GetUsersByRole implements UserService.
func (u *userService) GetUsersByRole(requesterRole string) ([]responses.UserResponse, error) {
	var roles []string

	switch requesterRole {
	case pkg.RoleAdmin:
		roles = []string{pkg.RoleKepalaPosyandu, pkg.RoleKader}
	case pkg.RoleKepalaPosyandu:
		roles = []string{pkg.RoleKader}
	default:
		return nil, errors.New("forbidden: role not allowed")
	}

	users, err := u.repo.FindUsersByRole(roles)
	if err != nil {
		return nil, err
	}

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
			CreatedBy:      v.CreatedBy,
			CreatedAt:      v.CreatedAt,
			UpdatedAt:      v.UpdatedAt,
		})
	}

	return result, nil
}

// DeleteUserByID implements UserService.
func (u *userService) DeleteUserByID(targetUserID int, updaterRole string) error {
	targetRole, err := u.repo.FindRoleById(targetUserID)

	if err != nil {
		return err
	}

	if !u.rolePermission(targetRole, updaterRole) {
		return fmt.Errorf("%s tidak boleh mengubah role ke %s", updaterRole, targetRole)
	}

	return u.repo.DeleteUser(targetUserID)
}

// CreateUser implements UserService.
func (u *userService) CreateUser(req requests.CreateUserRequest, createdBy string, locationID *int) (*responses.UserResponse, error) {

	if !u.rolePermission(createdBy, req.Role) {
		return nil, fmt.Errorf("role %s tidak diizinkan membuat user dengan role %s", createdBy, req.Role)
	}

	var url string
	var err error

	if req.ProfilePicture != nil {
		url, err = u.s3.UploadFile(req.ProfilePicture, "users")
		if err != nil {
			return nil, err
		}
	}

	hashedPassword, err := pkg.HashPassword(req.Password)

	if err != nil {
		return nil, err
	}

	var location int

	if createdBy == "admin" {
		if req.LocationID == nil {
			return nil, fmt.Errorf("admin harus menyertakan location_id")
		}
		location = *req.LocationID
	} else {
		if locationID == nil {
			return nil, fmt.Errorf("location_id tidak ditemukan di JWT")
		}
		location = *locationID
	}

	userMapping := models.User{
		LocationID:     location,
		Name:           req.Name,
		PhoneNumber:    req.PhoneNumber,
		Address:        req.Address,
		Nik:            req.Nik,
		ProfilePicture: url,
		Role:           req.Role,
		Password:       hashedPassword,
		CreatedBy:      createdBy,
	}

	user, err := u.repo.CreateUser(&userMapping)

	if err != nil {
		return nil, err
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
		CreatedBy:      user.CreatedBy,
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
	}

	return &response, nil
}

// DeleteUser implements UserService.
func (u *userService) DeleteCurrentUser(id int) error {
	return u.repo.DeleteUser(id)
}

// GetUser implements UserService.
func (u *userService) GetCurrentUser(id int) (*responses.UserResponse, error) {
	user, err := u.repo.GetUser(id)

	if err != nil {
		return nil, err
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
		CreatedBy:      user.CreatedBy,
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
	}

	return response, nil
}

// UpdateUser implements UserService.
func (u *userService) UpdateCurrentUser(id int, req requests.UpdateUserRequest) (*responses.UserResponse, error) {
	var url string
	var err error

	if req.ProfilePicture != nil {
		url, err = u.s3.UploadFile(req.ProfilePicture, "users")
		if err != nil {
			return nil, err
		}
	}

	userMapping := models.User{
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
		Address:     req.Address,
		Nik:         req.Nik,
	}

	if url != "" {
		userMapping.ProfilePicture = url
	}

	if req.Password != "" {
		hashedPassword, err := pkg.HashPassword(req.Password)
		if err != nil {
			return nil, err
		}
		userMapping.Password = hashedPassword
	}

	user, err := u.repo.UpdateUser(id, &userMapping)
	if err != nil {
		return nil, err
	}

	response := responses.UserResponse{
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
	}

	return &response, nil
}

// UpdateUserByID implements UserService.
func (u *userService) UpdateUserByID(
	targetUserID int,
	req requests.UpdateUserRequest,
	updaterRole string,
) (*responses.UserResponse, error) {

	role, err := u.repo.FindRoleById(targetUserID)

	if err != nil {
		return nil, err
	}

	if !u.rolePermission(updaterRole, role) {
		return nil, fmt.Errorf("%s tidak boleh mengubah role ke %s", updaterRole, req.Role)
	}

	var url string
	if req.ProfilePicture != nil {
		url, err = u.s3.UploadFile(req.ProfilePicture, "users")
		if err != nil {
			return nil, err
		}
	}

	userMapping := models.User{
		Name:        req.Name,
		LocationID:  *req.LocationID,
		PhoneNumber: req.PhoneNumber,
		Address:     req.Address,
		Nik:         req.Nik,
		Role:        req.Role,
	}

	if url != "" {
		userMapping.ProfilePicture = url
	}

	if req.Password != "" {
		hashedPassword, err := pkg.HashPassword(req.Password)
		if err != nil {
			return nil, err
		}
		userMapping.Password = hashedPassword
	}

	user, err := u.repo.UpdateUser(targetUserID, &userMapping)
	if err != nil {
		return nil, err
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

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo: repo}
}
