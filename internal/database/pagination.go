package database

type PaginationResponse struct {
	Data  interface{} `json:"data"`
	Total int         `json:"total"`
}
