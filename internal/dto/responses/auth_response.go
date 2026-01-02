package responses


type LoginResponse struct {
	Token TokenResponse `json:"token"`
}
