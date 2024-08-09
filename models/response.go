package models

// ResponseSuccess is a structure for successful responses.
type ResponseSuccess struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// ResponseError is a structure for failed responses.
type ResponseError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Errors  interface{} `json:"errors"`
}

// ValidationErrorDetail is a structure for validation error details.
type ValidationErrorDetail struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}
