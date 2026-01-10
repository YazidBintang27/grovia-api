package requests

type LoginRequest struct {
	PhoneNumber string `json:"phoneNumber" validate:"required,phone"`
	Password    string `json:"password" validate:"required,min=6"`
}

type ResetPasswordRequest struct {
	FirebaseToken   string `json:"firebaseToken"`
	PhoneNumber     string `json:"phoneNumber" validate:"required,phone"`
	Password        string `json:"password" validate:"required,min=6"`
	ConfirmPassword string `json:"confirmPassword" validate:"required,min=6"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}
