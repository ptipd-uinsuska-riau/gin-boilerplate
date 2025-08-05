package httplib

import (
	"github.com/gin-gonic/gin"
)

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

// SetSuccessResponse sets a successful response
func SetSuccessResponse(c *gin.Context, code int, httpStatus int, message string, data interface{}) {
	response := APIResponse{
		Success: true,
		Code:    code,
		Message: message,
		Data:    data,
	}
	c.JSON(httpStatus, response)
}

// SetErrorResponse sets an error response
func SetErrorResponse(c *gin.Context, code int, httpStatus int, message string) {
	response := APIErrorResponse{
		Success: false,
		Code:    code,
		Message: message,
	}
	c.JSON(httpStatus, response)
}

// SetErrorResponseWithError sets an error response with error details
func SetErrorResponseWithError(c *gin.Context, code int, httpStatus int, message string, err error) {
	response := APIErrorResponse{
		Success: false,
		Code:    code,
		Message: message,
		Error:   err.Error(),
	}
	c.JSON(httpStatus, response)
}

// SetPaginatedResponse sets a paginated response
func SetPaginatedResponse(c *gin.Context, code int, httpStatus int, message string, data interface{}, pagination PaginationMeta) {
	response := PaginatedResponse{
		Success:    true,
		Code:       code,
		Message:    message,
		Data:       data,
		Pagination: pagination,
	}
	c.JSON(httpStatus, response)
}

// CalculatePagination calculates pagination metadata
func CalculatePagination(page, limit int, total int64) PaginationMeta {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	return PaginationMeta{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}
}

// GetPaginationParams extracts pagination parameters from query
func GetPaginationParams(c *gin.Context) (page, limit int) {
	page = 1
	limit = 10

	if p := c.Query("page"); p != "" {
		if pageNum := parseInt(p); pageNum > 0 {
			page = pageNum
		}
	}

	if l := c.Query("limit"); l != "" {
		if limitNum := parseInt(l); limitNum > 0 && limitNum <= 100 {
			limit = limitNum
		}
	}

	return page, limit
}

// parseInt safely converts string to int
func parseInt(s string) int {
	var result int
	for _, char := range s {
		if char >= '0' && char <= '9' {
			result = result*10 + int(char-'0')
		} else {
			return 0
		}
	}
	return result
}
