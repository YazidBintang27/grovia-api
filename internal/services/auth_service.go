package services

import (
	"grovia/internal/dto/requests"
	"grovia/internal/dto/responses"
	"grovia/internal/repositories"
	"grovia/pkg"
)

type AuthService interface {
	Login(req requests.LoginRequest) (*responses.TokenResponse, error)
	ResetPassword(req requests.ResetPasswordRequest) error
	RefreshToken(refreshToken string) (*responses.TokenResponse, error)
}

type authService struct {
	repo repositories.AuthRepository
}

func (a *authService) RefreshToken(refreshToken string) (*responses.TokenResponse, error) {
	claims, err := pkg.ValidateToken(refreshToken)
	if err != nil {
		return nil, pkg.NewUnauthorizedError("Token tidak valid atau sudah kadaluarsa")
	}

	user, err := a.repo.FindByID(claims.UserID)
	if err != nil {
		return nil, pkg.NewNotFoundError("User tidak ditemukan")
	}

	accessToken, newRefreshToken, err := pkg.GenerateJWT(user.ID, user.LocationID, user.Role)
	if err != nil {
		return nil, pkg.NewInternalServerError("Gagal membuat token baru")
	}

	return &responses.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (a *authService) ResetPassword(req requests.ResetPasswordRequest) error {
	if err := pkg.ValidateStruct(req); err != nil {
		return pkg.NewBadRequestError(err.Error())
	}

	if req.Password != req.ConfirmPassword {
		return pkg.NewBadRequestError("Password dan Confirm Password tidak cocok")
	}

	_, err := VerifyFirebaseToken(req.FirebaseToken)
	if err != nil {
		return pkg.NewUnauthorizedError("Firebase token tidak valid")
	}

	hashed, err := pkg.HashPassword(req.Password)
	if err != nil {
		return pkg.NewInternalServerError("Gagal memproses password")
	}

	err = a.repo.ResetPassword(req.PhoneNumber, hashed)
	if err != nil {
		return pkg.NewInternalServerError("Gagal mereset password")
	}

	return nil
}

func (a *authService) Login(req requests.LoginRequest) (*responses.TokenResponse, error) {
	if err := pkg.ValidateStruct(req); err != nil {
		return nil, pkg.NewBadRequestError(err.Error())
	}

	user, err := a.repo.FindByPhoneNumber(req.PhoneNumber)
	if err != nil {
		return nil, pkg.NewNotFoundError("Nomor telepon tidak terdaftar")
	}

	if !pkg.CheckPassword(req.Password, user.Password) {
		return nil, pkg.NewUnauthorizedError("Password salah")
	}

	accessToken, refreshToken, err := pkg.GenerateJWT(user.ID, user.LocationID, user.Role)
	if err != nil {
		return nil, pkg.NewInternalServerError("Gagal membuat token")
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
