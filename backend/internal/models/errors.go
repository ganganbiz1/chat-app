package models

import "errors"

// Custom errors for the application
var (
	ErrRoomNotFound    = errors.New("room not found")
	ErrMessageNotFound = errors.New("message not found")
	ErrInvalidInput    = errors.New("invalid input")
	ErrDatabaseError   = errors.New("database error")
	ErrCacheError      = errors.New("cache error")
	ErrValidationError = errors.New("validation error")
)

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Code    int    `json:"code,omitempty"`
}

// SuccessResponse represents a success response with data
type SuccessResponse struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message,omitempty"`
}

// PaginationResponse represents a paginated response
type PaginationResponse struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	Total      int64       `json:"total"`
	TotalPages int         `json:"total_pages"`
}
