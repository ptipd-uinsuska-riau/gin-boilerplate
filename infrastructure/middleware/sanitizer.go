package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"reflect"
	"strings"

	"gin-boilerplate/infrastructure/errors"
	"gin-boilerplate/infrastructure/sanitizer"

	"github.com/gin-gonic/gin"
)

// SanitizerConfig holds configuration for the sanitizer middleware
type SanitizerConfig struct {
	// Enable XSS detection and blocking
	EnableXSSDetection bool
	// Enable SQL injection detection and blocking
	EnableSQLInjectionDetection bool
	// Maximum length for string fields
	MaxStringLength int
	// Fields to skip sanitization (case-insensitive)
	SkipFields []string
	// Enable strict mode (block requests with detected threats)
	StrictMode bool
}

// DefaultSanitizerConfig returns default configuration
func DefaultSanitizerConfig() SanitizerConfig {
	return SanitizerConfig{
		EnableXSSDetection:          true,
		EnableSQLInjectionDetection: true,
		MaxStringLength:             1000,
		SkipFields:                  []string{"password", "token", "secret"},
		StrictMode:                  true,
	}
}

// SanitizerMiddleware creates a middleware that sanitizes request body
func SanitizerMiddleware(config ...SanitizerConfig) gin.HandlerFunc {
	cfg := DefaultSanitizerConfig()
	if len(config) > 0 {
		cfg = config[0]
	}

	sanitizerInstance := sanitizer.NewSanitizer()

	return func(c *gin.Context) {
		// Only process JSON requests with body
		if c.Request.Method == "GET" || c.Request.Method == "DELETE" ||
			c.Request.ContentLength == 0 ||
			!strings.Contains(c.GetHeader("Content-Type"), "application/json") {
			c.Next()
			return
		}

		// Read the request body
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			errors.AbortWithError(c, errors.NewBadRequestError("Failed to read request body"))
			return
		}

		// Close the original body
		c.Request.Body.Close()

		// Parse JSON
		var jsonData interface{}
		if err := json.Unmarshal(body, &jsonData); err != nil {
			// If it's not valid JSON, restore the original body and continue
			c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
			c.Next()
			return
		}

		// Sanitize the JSON data
		sanitizedData, blocked := sanitizeJSONData(jsonData, sanitizerInstance, cfg)
		if blocked && cfg.StrictMode {
			errors.AbortWithError(c, errors.NewBadRequestError("Request contains potentially malicious content"))
			return
		}

		// Marshal the sanitized data back to JSON
		sanitizedBody, err := json.Marshal(sanitizedData)
		if err != nil {
			errors.AbortWithError(c, errors.NewInternalServerError("Failed to process request", err))
			return
		}

		// Replace the request body with sanitized data
		c.Request.Body = io.NopCloser(bytes.NewBuffer(sanitizedBody))
		c.Request.ContentLength = int64(len(sanitizedBody))

		c.Next()
	}
}

// sanitizeJSONData recursively sanitizes JSON data
func sanitizeJSONData(data interface{}, sanitizerInstance *sanitizer.Sanitizer, config SanitizerConfig) (interface{}, bool) {
	blocked := false

	switch v := data.(type) {
	case map[string]interface{}:
		result := make(map[string]interface{})
		for key, value := range v {
			// Skip certain fields
			if shouldSkipField(key, config.SkipFields) {
				result[key] = value
				continue
			}

			sanitizedValue, wasBlocked := sanitizeJSONData(value, sanitizerInstance, config)
			if wasBlocked {
				blocked = true
			}
			result[key] = sanitizedValue
		}
		return result, blocked

	case []interface{}:
		result := make([]interface{}, len(v))
		for i, item := range v {
			sanitizedItem, wasBlocked := sanitizeJSONData(item, sanitizerInstance, config)
			if wasBlocked {
				blocked = true
			}
			result[i] = sanitizedItem
		}
		return result, blocked

	case string:
		return sanitizeStringValue(v, sanitizerInstance, config)

	default:
		return data, false
	}
}

// sanitizeStringValue sanitizes a string value and checks for threats
func sanitizeStringValue(value string, sanitizerInstance *sanitizer.Sanitizer, config SanitizerConfig) (string, bool) {
	blocked := false

	// Check length
	if config.MaxStringLength > 0 && len(value) > config.MaxStringLength {
		blocked = true
		value = value[:config.MaxStringLength]
	}

	// Check for XSS
	if config.EnableXSSDetection && sanitizerInstance.DetectXSS(value) {
		blocked = true
	}

	// Check for SQL injection
	if config.EnableSQLInjectionDetection && sanitizerInstance.DetectSQLInjection(value) {
		blocked = true
	}

	// Sanitize the string
	sanitizedValue := sanitizerInstance.SanitizeString(value)

	return sanitizedValue, blocked
}

// shouldSkipField checks if a field should be skipped during sanitization
func shouldSkipField(fieldName string, skipFields []string) bool {
	fieldNameLower := strings.ToLower(fieldName)
	for _, skipField := range skipFields {
		if strings.ToLower(skipField) == fieldNameLower {
			return true
		}
	}
	return false
}

// SanitizeStructFields sanitizes string fields in a struct using reflection
func SanitizeStructFields(data interface{}, config ...SanitizerConfig) interface{} {
	cfg := DefaultSanitizerConfig()
	if len(config) > 0 {
		cfg = config[0]
	}

	sanitizerInstance := sanitizer.NewSanitizer()

	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return data
	}

	// Create a copy of the struct
	result := reflect.New(v.Type()).Elem()
	result.Set(v)

	// Iterate through fields
	for i := 0; i < v.NumField(); i++ {
		field := result.Field(i)
		fieldType := v.Type().Field(i)

		// Skip unexported fields
		if !field.CanSet() {
			continue
		}

		// Skip certain fields
		if shouldSkipField(fieldType.Name, cfg.SkipFields) {
			continue
		}

		// Only process string fields
		if field.Kind() == reflect.String {
			originalValue := field.String()
			sanitizedValue := sanitizerInstance.SanitizeString(originalValue)
			field.SetString(sanitizedValue)
		}
	}

	return result.Interface()
}

// HeaderSanitizerMiddleware sanitizes HTTP headers
func HeaderSanitizerMiddleware() gin.HandlerFunc {
	sanitizerInstance := sanitizer.NewSanitizer()

	return func(c *gin.Context) {
		// Sanitize common headers that might contain user input
		headersToSanitize := []string{
			"User-Agent",
			"Referer",
			"X-Forwarded-For",
			"X-Real-IP",
		}

		for _, headerName := range headersToSanitize {
			headerValue := c.GetHeader(headerName)
			if headerValue != "" {
				sanitizedValue := sanitizerInstance.SanitizeString(headerValue)
				c.Request.Header.Set(headerName, sanitizedValue)
			}
		}

		c.Next()
	}
}

// QueryParamSanitizerMiddleware sanitizes query parameters
func QueryParamSanitizerMiddleware() gin.HandlerFunc {
	sanitizerInstance := sanitizer.NewSanitizer()

	return func(c *gin.Context) {
		// Get all query parameters
		queryParams := c.Request.URL.Query()

		// Sanitize each parameter
		for key, values := range queryParams {
			for i, value := range values {
				sanitizedValue := sanitizerInstance.SanitizeString(value)
				queryParams[key][i] = sanitizedValue
			}
		}

		// Update the request URL with sanitized parameters
		c.Request.URL.RawQuery = queryParams.Encode()

		c.Next()
	}
}
