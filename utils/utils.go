package utils

import (
	"strconv"
	"strings"
	"time"
)

// StringUnitToDuration converts string with unit to time.Duration
func StringUnitToDuration(s string) time.Duration {
	if s == "" {
		return time.Minute
	}

	// Parse duration string (e.g., "1m", "30s", "1h")
	duration, err := time.ParseDuration(s)
	if err != nil {
		// If parsing fails, try to extract number and assume minutes
		if num, err := strconv.Atoi(strings.TrimSpace(s)); err == nil {
			return time.Duration(num) * time.Minute
		}
		// Default to 1 minute if all parsing fails
		return time.Minute
	}

	return duration
}

// Contains checks if a string slice contains a specific string
func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// RemoveEmpty removes empty strings from a slice
func RemoveEmpty(slice []string) []string {
	var result []string
	for _, s := range slice {
		if strings.TrimSpace(s) != "" {
			result = append(result, s)
		}
	}
	return result
}

// TruncateString truncates a string to a maximum length
func TruncateString(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}
	return s[:maxLength] + "..."
}

// IsValidEmail performs basic email validation
func IsValidEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

// GetEnvOrDefault returns environment variable value or default if not set
func GetEnvOrDefault(key, defaultValue string) string {
	if value := strings.TrimSpace(key); value != "" {
		return value
	}
	return defaultValue
}

// ErrorLogFormat is a constant for consistent error logging format
const ErrorLogFormat = "Error: %s | Context: %s | Function: %s"
