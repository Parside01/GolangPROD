package models

type ErrorResponse struct {
	Err string `json:"err"`
}

type StatusOKResponse struct {
	Status string `json:"status"`
}
