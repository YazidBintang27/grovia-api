package services

import (
	"errors"
	"grovia/internal/dto/requests"
	"grovia/internal/dto/responses"
	"grovia/internal/repositories"
	"grovia/pkg"
)

type AuthService interface {
	Login(req requests.LoginRequest) (*responses.TokenResponse, error)
	ResetPassword(req requests.ResetPasswordRequest) error
}

type authService struct {
	repo repositories.AuthRepository
}

// ResetPassword implements AuthService.
func (a *authService) ResetPassword(req requests.ResetPasswordRequest) error {
	if req.Password != req.ConfirmPassword {
		return errors.New("password and confirm password do not match")
	}

	_, err := VerifyFirebaseToken(req.FirebaseToken)

	if err != nil {
		return err
	}

	hashed, err := pkg.HashPassword(req.Password)

	if err != nil {
		return err
	}

	err = a.repo.ResetPassword(req.PhoneNumber, hashed)

	if err != nil {
		return err
	}

	return nil
}

// Login implements AuthService.
func (a *authService) Login(req requests.LoginRequest) (*responses.TokenResponse, error) {
	user, err := a.repo.FindByPhoneNumber(req.PhoneNumber)

	if err != nil {
		return nil, errors.New("invalid Phone Number or Password")
	}

	if !pkg.CheckPassword(req.Password, user.Password) {
		return nil, errors.New("invalid Password")
	}

	accessToken, refreshToken, err := pkg.GenerateJWT(user.ID, user.Role)

	if err != nil {
		return nil, errors.New("invalid Token")
	}

	tokenResponse := responses.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return &tokenResponse, nil
}

func NewAuthService(repo repositories.AuthRepository) AuthService {
	return &authService{repo: repo}
}
