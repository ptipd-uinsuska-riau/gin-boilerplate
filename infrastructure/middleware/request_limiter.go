package middleware

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"gin-boilerplate/infrastructure/errors"

	"github.com/gin-gonic/gin"
)

// RequestLimiterConfig holds configuration for request limiting
type RequestLimiterConfig struct {
	// Maximum request body size in bytes
	MaxBodySize int64
	// Maximum number of form fields
	MaxFormFields int
	// Maximum form memory size in bytes
	MaxFormMemory int64
	// Maximum number of multipart form files
	MaxFormFiles int
	// Maximum header size in bytes
	MaxHeaderSize int64
	// Maximum URL length
	MaxURLLength int
	// Maximum query parameter count
	MaxQueryParams int
	// Skip size check for certain paths
	SkipPaths []string
}

// DefaultRequestLimiterConfig returns default configuration
func DefaultRequestLimiterConfig() RequestLimiterConfig {
	return RequestLimiterConfig{
		MaxBodySize:    10 << 20, // 10 MB
		MaxFormFields:  100,
		MaxFormMemory:  32 << 20, // 32 MB
		MaxFormFiles:   10,
		MaxHeaderSize:  1 << 20, // 1 MB
		MaxURLLength:   2048,
		MaxQueryParams: 50,
		SkipPaths:      []string{"/api/v1/health"},
	}
}

// RequestLimiterMiddleware creates a middleware that limits request size and complexity
func RequestLimiterMiddleware(config ...RequestLimiterConfig) gin.HandlerFunc {
	cfg := DefaultRequestLimiterConfig()
	if len(config) > 0 {
		cfg = config[0]
	}

	return func(c *gin.Context) {
		// Check if path should be skipped
		if shouldSkipPath(c.Request.URL.Path, cfg.SkipPaths) {
			c.Next()
			return
		}

		// Check URL length
		if cfg.MaxURLLength > 0 && len(c.Request.URL.String()) > cfg.MaxURLLength {
			errors.AbortWithError(c, errors.NewBadRequestError(
				fmt.Sprintf("URL too long. Maximum allowed: %d characters", cfg.MaxURLLength)))
			return
		}

		// Check query parameter count
		if cfg.MaxQueryParams > 0 && len(c.Request.URL.Query()) > cfg.MaxQueryParams {
			errors.AbortWithError(c, errors.NewBadRequestError(
				fmt.Sprintf("Too many query parameters. Maximum allowed: %d", cfg.MaxQueryParams)))
			return
		}

		// Check header size
		if cfg.MaxHeaderSize > 0 {
			headerSize := calculateHeaderSize(c.Request)
			if headerSize > cfg.MaxHeaderSize {
				errors.AbortWithError(c, errors.NewBadRequestError(
					fmt.Sprintf("Headers too large. Maximum allowed: %d bytes", cfg.MaxHeaderSize)))
				return
			}
		}

		// Check content length
		if c.Request.ContentLength > cfg.MaxBodySize {
			errors.AbortWithError(c, errors.NewBadRequestError(
				fmt.Sprintf("Request body too large. Maximum allowed: %d bytes", cfg.MaxBodySize)))
			return
		}

		// Wrap the request body with a limited reader
		if c.Request.Body != nil {
			c.Request.Body = &limitedReadCloser{
				ReadCloser: c.Request.Body,
				limit:      cfg.MaxBodySize,
				read:       0,
			}
		}

		// Set max form memory for multipart forms
		if cfg.MaxFormMemory > 0 {
			c.Request.ParseMultipartForm(cfg.MaxFormMemory)
		}

		// Check form complexity for POST/PUT requests
		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "PATCH" {
			if err := checkFormComplexity(c, cfg); err != nil {
				errors.AbortWithError(c, err)
				return
			}
		}

		c.Next()
	}
}

// limitedReadCloser wraps an io.ReadCloser with size limiting
type limitedReadCloser struct {
	io.ReadCloser
	limit int64
	read  int64
}

func (lrc *limitedReadCloser) Read(p []byte) (int, error) {
	if lrc.read >= lrc.limit {
		return 0, fmt.Errorf("request body size limit exceeded")
	}

	// Limit the read to not exceed the total limit
	if int64(len(p)) > lrc.limit-lrc.read {
		p = p[:lrc.limit-lrc.read]
	}

	n, err := lrc.ReadCloser.Read(p)
	lrc.read += int64(n)

	if lrc.read > lrc.limit {
		return n, fmt.Errorf("request body size limit exceeded")
	}

	return n, err
}

// calculateHeaderSize calculates the total size of HTTP headers
func calculateHeaderSize(req *http.Request) int64 {
	size := int64(0)
	
	// Add request line size
	size += int64(len(req.Method) + len(req.URL.String()) + len(req.Proto) + 4) // spaces and CRLF
	
	// Add headers size
	for name, values := range req.Header {
		for _, value := range values {
			size += int64(len(name) + len(value) + 4) // ": " and CRLF
		}
	}
	
	// Add final CRLF
	size += 2
	
	return size
}

// checkFormComplexity checks form field and file limits
func checkFormComplexity(c *gin.Context, cfg RequestLimiterConfig) *errors.AppError {
	contentType := c.GetHeader("Content-Type")
	
	// Check multipart forms
	if strings.Contains(contentType, "multipart/form-data") {
		if err := c.Request.ParseMultipartForm(cfg.MaxFormMemory); err != nil {
			return errors.NewBadRequestError("Failed to parse multipart form: " + err.Error())
		}

		// Check number of form fields
		if cfg.MaxFormFields > 0 && c.Request.MultipartForm != nil {
			fieldCount := len(c.Request.MultipartForm.Value)
			if fieldCount > cfg.MaxFormFields {
				return errors.NewBadRequestError(
					fmt.Sprintf("Too many form fields. Maximum allowed: %d", cfg.MaxFormFields))
			}
		}

		// Check number of files
		if cfg.MaxFormFiles > 0 && c.Request.MultipartForm != nil {
			fileCount := len(c.Request.MultipartForm.File)
			if fileCount > cfg.MaxFormFiles {
				return errors.NewBadRequestError(
					fmt.Sprintf("Too many form files. Maximum allowed: %d", cfg.MaxFormFiles))
			}
		}
	}

	// Check URL-encoded forms
	if strings.Contains(contentType, "application/x-www-form-urlencoded") {
		if err := c.Request.ParseForm(); err != nil {
			return errors.NewBadRequestError("Failed to parse form: " + err.Error())
		}

		// Check number of form fields
		if cfg.MaxFormFields > 0 {
			fieldCount := len(c.Request.Form)
			if fieldCount > cfg.MaxFormFields {
				return errors.NewBadRequestError(
					fmt.Sprintf("Too many form fields. Maximum allowed: %d", cfg.MaxFormFields))
			}
		}
	}

	return nil
}

// shouldSkipPath checks if a path should be skipped
func shouldSkipPath(path string, skipPaths []string) bool {
	for _, skipPath := range skipPaths {
		if strings.HasPrefix(path, skipPath) {
			return true
		}
	}
	return false
}

// FileSizeLimiterMiddleware creates a middleware specifically for file upload size limiting
func FileSizeLimiterMiddleware(maxFileSize int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		contentType := c.GetHeader("Content-Type")
		
		if strings.Contains(contentType, "multipart/form-data") {
			// Parse the multipart form with size limit
			if err := c.Request.ParseMultipartForm(maxFileSize); err != nil {
				errors.AbortWithError(c, errors.NewBadRequestError(
					fmt.Sprintf("File too large. Maximum allowed: %d bytes", maxFileSize)))
				return
			}

			// Check individual file sizes
			if c.Request.MultipartForm != nil {
				for _, fileHeaders := range c.Request.MultipartForm.File {
					for _, fileHeader := range fileHeaders {
						if fileHeader.Size > maxFileSize {
							errors.AbortWithError(c, errors.NewBadRequestError(
								fmt.Sprintf("File '%s' too large. Maximum allowed: %d bytes", 
									fileHeader.Filename, maxFileSize)))
							return
						}
					}
				}
			}
		}

		c.Next()
	}
}

// ResponseSizeLimiterMiddleware limits response size
func ResponseSizeLimiterMiddleware(maxResponseSize int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Wrap the response writer
		lrw := &limitedResponseWriter{
			ResponseWriter: c.Writer,
			limit:          maxResponseSize,
			written:        0,
		}
		c.Writer = lrw

		c.Next()
	}
}

// limitedResponseWriter wraps gin.ResponseWriter with size limiting
type limitedResponseWriter struct {
	gin.ResponseWriter
	limit   int64
	written int64
}

func (lrw *limitedResponseWriter) Write(data []byte) (int, error) {
	if lrw.written+int64(len(data)) > lrw.limit {
		// Write as much as possible within the limit
		remaining := lrw.limit - lrw.written
		if remaining > 0 {
			n, err := lrw.ResponseWriter.Write(data[:remaining])
			lrw.written += int64(n)
			return n, err
		}
		return 0, fmt.Errorf("response size limit exceeded")
	}

	n, err := lrw.ResponseWriter.Write(data)
	lrw.written += int64(n)
	return n, err
}

func (lrw *limitedResponseWriter) WriteString(s string) (int, error) {
	return lrw.Write([]byte(s))
}

// RequestTimeoutMiddleware adds request timeout
func RequestTimeoutMiddleware(timeoutHeader string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if timeoutStr := c.GetHeader(timeoutHeader); timeoutStr != "" {
			if timeout, err := strconv.Atoi(timeoutStr); err == nil && timeout > 0 {
				// Set a reasonable maximum timeout (e.g., 300 seconds)
				if timeout > 300 {
					timeout = 300
				}
				c.Header("X-Request-Timeout", strconv.Itoa(timeout))
			}
		}
		c.Next()
	}
}
