package models

type CommonResponse struct {
	Code        int    `json:"code"`
	Message     string `json:"message"`
	Error       string `json:"error,omitempty"`
	ErrorDetail string `json:"error_detail,omitempty"`
}
