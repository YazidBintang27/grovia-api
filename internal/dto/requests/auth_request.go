package requests

type LoginRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
	Password    string `json:"password" binding:"required"`
}

type ResetPasswordRequest struct {
	FirebaseToken   string `json:"firebaseToken"`
	PhoneNumber     string `json:"phoneNumber"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}
