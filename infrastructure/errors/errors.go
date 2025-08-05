package errors

import (
	"fmt"
	"net/http"
)

// AppError represents an application error with HTTP status code
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// NewAppError creates a new application error
func NewAppError(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// Predefined error constructors
func NewBadRequestError(message string) *AppError {
	return NewAppError(http.StatusBadRequest, message, nil)
}

func NewUnauthorizedError(message string) *AppError {
	return NewAppError(http.StatusUnauthorized, message, nil)
}

func NewForbiddenError(message string) *AppError {
	return NewAppError(http.StatusForbidden, message, nil)
}

func NewNotFoundError(message string) *AppError {
	return NewAppError(http.StatusNotFound, message, nil)
}

func NewConflictError(message string) *AppError {
	return NewAppError(http.StatusConflict, message, nil)
}

func NewInternalServerError(message string, err error) *AppError {
	return NewAppError(http.StatusInternalServerError, message, err)
}

func NewValidationError(message string) *AppError {
	return NewAppError(http.StatusUnprocessableEntity, message, nil)
}

func NewRateLimitError(message string) *AppError {
	return NewAppError(http.StatusTooManyRequests, message, nil)
}

// Common error messages
const (
	ErrInvalidInput     = "Invalid input provided"
	ErrUnauthorized     = "Unauthorized access"
	ErrForbidden        = "Access forbidden"
	ErrNotFound         = "Resource not found"
	ErrConflict         = "Resource already exists"
	ErrInternalServer   = "Internal server error"
	ErrValidation       = "Validation failed"
	ErrRateLimit        = "Rate limit exceeded"
	ErrInvalidToken     = "Invalid or expired token"
	ErrInvalidPassword  = "Invalid password"
	ErrUserNotFound     = "User not found"
	ErrEmailExists      = "Email already exists"
)
