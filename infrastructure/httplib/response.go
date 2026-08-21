package httplib

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)


// catatKegagalanServer mencatat kegagalan 5xx yang TIDAK melewati ErrorHandler.
//
// Ada dua jalur menjawab error di boilerplate ini: AbortWithError (menaruh error
// di c.Errors, lalu ErrorHandler mencatatnya) dan pemanggilan SetErrorResponse
// langsung dari handler. Jalur kedua TIDAK PERNAH tercatat di mana pun —
// kegagalan 5xx, termasuk kegagalan basis data, hanya terkirim ke pemanggil lalu
// hilang. Beberapa aplikasi turunan memakai jalur kedua secara eksklusif,
// sehingga di sana tidak ada satu pun jejak kegagalan server.
//
// Penjaga len(c.Errors) == 0 memisahkan keduanya: bila error sudah ditaruh di
// c.Errors berarti ErrorHandler yang mencatatnya, dan mencatat lagi di sini
// hanya melipatgandakan baris log — persis yang membuat log sulit dibaca saat
// gangguan sedang ditelusuri.
//
// 4xx sengaja TIDAK dicatat: itu kesalahan pemanggil, dan mencatatnya membuat
// UUID salah ketik ikut menjadi baris ERROR.
//
// Pesannya saja yang dicatat, bukan detail error mentah: pesan mentah dapat
// memuat isi baris basis data. Itu urusan jejak SQL, bukan jejak permintaan.
func catatKegagalanServer(c *gin.Context, httpStatus int, message string) {
	if httpStatus < 500 || len(c.Errors) > 0 {
		return
	}
	logrus.WithFields(logrus.Fields{
		"method": c.Request.Method,
		"path":   c.Request.URL.Path,
		"status": httpStatus,
	}).Error("Kegagalan server: " + message)
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
	catatKegagalanServer(c, httpStatus, message)
	response := APIErrorResponse{
		Success: false,
		Code:    code,
		Message: message,
	}
	c.JSON(httpStatus, response)
}

// SetErrorResponseWithError sets an error response with error details
func SetErrorResponseWithError(c *gin.Context, code int, httpStatus int, message string, err error) {
	catatKegagalanServer(c, httpStatus, message)
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
