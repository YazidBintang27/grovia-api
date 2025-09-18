package responses

type BaseResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data"`
	Error   any    `json:"error"`
}
