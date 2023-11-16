package models

type ErrorResponse struct {
	Details string `json:"details"`
	Error   string `json:"error"`
}
