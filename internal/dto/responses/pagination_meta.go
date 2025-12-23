package responses

type PaginationMeta struct {
	Page int `json:"page"`
	Limit int `json:"limit"`
	TotalData int `json:"totalData"`
	TotalPage int `json:"totalPage"`
}