package models

type PaginationMillionaireDto struct {
	Millionaires []Millionaire `json:"millionaires"`
	Total        int           `json:"total"`
	Page         int           `json:"page"`
	PageSize     int           `json:"pageSize"`
}
