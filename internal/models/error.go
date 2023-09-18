package models

type ErrorResponse struct {
	Message string                 `json:"message"`
	Errors  []ErrorMessageResponse `json:"errors"`
}

type ErrorMessageResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
