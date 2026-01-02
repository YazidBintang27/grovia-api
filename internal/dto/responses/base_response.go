package responses

type BaseResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Data    any             `json:"data,omitempty"`
	Meta    *PaginationMeta `json:"meta,omitempty"`
	Error   any             `json:"error,omitempty"`
}
