package dto

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type Response[T any] struct {
	Data T `json:"data"`
}

type Page[T any] struct {
	Data  []T   `json:"data"`
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
	Total int64 `json:"total"`
}

type PageQuery struct {
	Page  int `query:"page"  validate:"omitempty,min=1"`
	Limit int `query:"limit" validate:"omitempty,min=1,max=100"`
}