package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"gin-boilerplate/infrastructure/config"
	"gin-boilerplate/infrastructure/errors"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSanitizerMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    string
		expectedStatus int
		strictMode     bool
		shouldBlock    bool
	}{
		{
			name:           "Normal request should pass",
			requestBody:    `{"name": "John Doe", "email": "john@example.com"}`,
			expectedStatus: http.StatusOK,
			strictMode:     true,
			shouldBlock:    false,
		},
		{
			name:           "XSS attempt should be blocked in strict mode",
			requestBody:    `{"name": "<script>alert('xss')</script>", "email": "john@example.com"}`,
			expectedStatus: http.StatusBadRequest,
			strictMode:     true,
			shouldBlock:    true,
		},
		{
			name:           "XSS attempt should be sanitized in non-strict mode",
			requestBody:    `{"name": "<script>alert('xss')</script>", "email": "john@example.com"}`,
			expectedStatus: http.StatusOK,
			strictMode:     false,
			shouldBlock:    false,
		},
		{
			name:           "SQL injection attempt should be blocked",
			requestBody:    `{"name": "John'; DROP TABLE users; --", "email": "john@example.com"}`,
			expectedStatus: http.StatusBadRequest,
			strictMode:     true,
			shouldBlock:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			r.Use(errors.ErrorHandler())

			config := SanitizerConfig{
				EnableXSSDetection:          true,
				EnableSQLInjectionDetection: true,
				MaxStringLength:             1000,
				StrictMode:                  tt.strictMode,
				SkipFields:                  []string{"password"},
			}

			r.Use(SanitizerMiddleware(config))
			r.POST("/test", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"status": "ok"})
			})

			req, _ := http.NewRequest("POST", "/test", bytes.NewBufferString(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestRequestLimiterMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    string
		contentLength  int64
		expectedStatus int
		maxBodySize    int64
	}{
		{
			name:           "Normal request should pass",
			requestBody:    `{"name": "John Doe"}`,
			expectedStatus: http.StatusOK,
			maxBodySize:    1024,
		},
		{
			name:           "Large request should be rejected",
			requestBody:    strings.Repeat("a", 2000),
			expectedStatus: http.StatusBadRequest,
			maxBodySize:    1024,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			r.Use(errors.ErrorHandler())

			config := RequestLimiterConfig{
				MaxBodySize:    tt.maxBodySize,
				MaxFormFields:  10,
				MaxFormMemory:  1024,
				MaxFormFiles:   5,
				MaxHeaderSize:  1024,
				MaxURLLength:   2048,
				MaxQueryParams: 10,
			}

			r.Use(RequestLimiterMiddleware(config))
			r.POST("/test", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"status": "ok"})
			})

			req, _ := http.NewRequest("POST", "/test", bytes.NewBufferString(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			req.ContentLength = int64(len(tt.requestBody))

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestHeaderSanitizerMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(HeaderSanitizerMiddleware())
	r.GET("/test", func(c *gin.Context) {
		userAgent := c.GetHeader("User-Agent")
		c.JSON(http.StatusOK, gin.H{"user_agent": userAgent})
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("User-Agent", "<script>alert('xss')</script>Mozilla/5.0")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	// The response should contain sanitized user agent
	assert.Contains(t, w.Body.String(), "Mozilla/5.0")
	assert.NotContains(t, w.Body.String(), "<script>")
}

func TestQueryParamSanitizerMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(QueryParamSanitizerMiddleware())
	r.GET("/test", func(c *gin.Context) {
		search := c.Query("search")
		c.JSON(http.StatusOK, gin.H{"search": search})
	})

	req, _ := http.NewRequest("GET", "/test?search=<script>alert('xss')</script>test", nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	// The response should contain sanitized search parameter
	assert.Contains(t, w.Body.String(), "test")
	assert.NotContains(t, w.Body.String(), "<script>")
}

func TestFileSizeLimiterMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		maxFileSize    int64
		expectedStatus int
	}{
		{
			name:           "Small file should pass",
			maxFileSize:    1024 * 1024, // 1MB
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			r.Use(FileSizeLimiterMiddleware(tt.maxFileSize))
			r.POST("/upload", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"status": "ok"})
			})

			// Create a simple form data request
			body := bytes.NewBufferString("--boundary\r\nContent-Disposition: form-data; name=\"file\"; filename=\"test.txt\"\r\nContent-Type: text/plain\r\n\r\ntest content\r\n--boundary--\r\n")
			req, _ := http.NewRequest("POST", "/upload", body)
			req.Header.Set("Content-Type", "multipart/form-data; boundary=boundary")

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestSecurityConfigIntegration(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Test with actual config structure
	config.Conf.Security = config.SecurityConfig{
		EnableSanitization:          true,
		EnableXSSDetection:          true,
		EnableSQLInjectionDetection: true,
		MaxStringLength:             100,
		StrictMode:                  true,
		MaxBodySize:                 1024,
		MaxFormFields:               10,
		MaxFormMemory:               1024,
		MaxFormFiles:                5,
		MaxHeaderSize:               1024,
		MaxURLLength:                2048,
		MaxQueryParams:              10,
		MaxFileSize:                 1024 * 1024,
		MaxResponseSize:             1024 * 1024,
	}

	r := gin.New()
	r.Use(errors.ErrorHandler())

	// Apply security middlewares like in the router
	sanitizerConfig := SanitizerConfig{
		EnableXSSDetection:          config.Conf.Security.EnableXSSDetection,
		EnableSQLInjectionDetection: config.Conf.Security.EnableSQLInjectionDetection,
		MaxStringLength:             config.Conf.Security.MaxStringLength,
		StrictMode:                  config.Conf.Security.StrictMode,
		SkipFields:                  []string{"password", "token", "secret"},
	}
	r.Use(SanitizerMiddleware(sanitizerConfig))

	requestLimiterConfig := RequestLimiterConfig{
		MaxBodySize:    config.Conf.Security.MaxBodySize,
		MaxFormFields:  config.Conf.Security.MaxFormFields,
		MaxFormMemory:  config.Conf.Security.MaxFormMemory,
		MaxFormFiles:   config.Conf.Security.MaxFormFiles,
		MaxHeaderSize:  config.Conf.Security.MaxHeaderSize,
		MaxURLLength:   config.Conf.Security.MaxURLLength,
		MaxQueryParams: config.Conf.Security.MaxQueryParams,
		SkipPaths:      []string{"/health"},
	}
	r.Use(RequestLimiterMiddleware(requestLimiterConfig))

	r.POST("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Test malicious request
	maliciousBody := `{"name": "<script>alert('xss')</script>", "description": "test"}`
	req, _ := http.NewRequest("POST", "/test", bytes.NewBufferString(maliciousBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Should be blocked due to XSS detection in strict mode
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
