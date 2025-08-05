package primitive

import (
	"time"

	"github.com/google/uuid"
)

// BaseModel contains common fields for all models
type BaseModel struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// APIResponse represents a standard API response
type APIResponse struct {
	Success bool        `json:"success"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// APIErrorResponse represents a standard API error response
type APIErrorResponse struct {
	Success bool   `json:"success"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// PaginationMeta represents pagination metadata
type PaginationMeta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// PaginatedResponse represents a paginated API response
type PaginatedResponse struct {
	Success    bool           `json:"success"`
	Code       int            `json:"code"`
	Message    string         `json:"message"`
	Data       interface{}    `json:"data"`
	Pagination PaginationMeta `json:"pagination"`
}

// HealthResponse represents health check response
type HealthResponse struct {
	Database string `json:"database"`
	Redis    string `json:"redis"`
	Uptime   string `json:"uptime,omitempty"`
}

// LoginRequest represents login request
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// LoginResponse represents login response
type LoginResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	User      UserInfo  `json:"user"`
}

// RegisterRequest represents registration request
type RegisterRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=100"`
}

// UserInfo represents user information
type UserInfo struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
}

// UpdateUserRequest represents update user request
type UpdateUserRequest struct {
	Name  string `json:"name" validate:"omitempty,min=2,max=100"`
	Email string `json:"email" validate:"omitempty,email"`
}

// ChangePasswordRequest represents change password request
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=6,max=100"`
}

// QueryParams represents common query parameters
type QueryParams struct {
	Page   int    `form:"page" validate:"omitempty,gte=1"`
	Limit  int    `form:"limit" validate:"omitempty,gte=1,lte=100"`
	Search string `form:"search" validate:"omitempty,max=100"`
	Sort   string `form:"sort" validate:"omitempty,oneof=asc desc"`
	SortBy string `form:"sort_by" validate:"omitempty,max=50"`
}

// IDParam represents ID parameter
type IDParam struct {
	ID uuid.UUID `uri:"id" validate:"required,uuid"`
}
