package sanitizer

import (
	"html"
	"regexp"
	"strings"
	"unicode"
)

// Sanitizer provides methods for sanitizing user input
type Sanitizer struct {
	// HTML tag removal regex
	htmlTagRegex *regexp.Regexp
	// SQL injection patterns
	sqlInjectionRegex *regexp.Regexp
	// Script tag removal regex
	scriptTagRegex *regexp.Regexp
	// Dangerous characters regex
	dangerousCharsRegex *regexp.Regexp
}

// NewSanitizer creates a new sanitizer instance
func NewSanitizer() *Sanitizer {
	return &Sanitizer{
		htmlTagRegex:        regexp.MustCompile(`<[^>]*>`),
		sqlInjectionRegex:   regexp.MustCompile(`(?i)(union\s+select|select\s+.*\s+from|\'\s*or\s+\'\d+\'\s*=\s*\'\d+|drop\s+table|insert\s+into|update\s+.*\s+set|delete\s+from|exec\s*\(|execute\s*\()`),
		scriptTagRegex:      regexp.MustCompile(`(?i)<script[^>]*>.*?</script>`),
		dangerousCharsRegex: regexp.MustCompile(`[<>\"'&]`),
	}
}

// SanitizeString sanitizes a string by removing dangerous characters and HTML tags
func (s *Sanitizer) SanitizeString(input string) string {
	if input == "" {
		return input
	}

	// Remove script tags
	sanitized := s.scriptTagRegex.ReplaceAllString(input, "")

	// Remove HTML tags
	sanitized = s.htmlTagRegex.ReplaceAllString(sanitized, "")

	// HTML escape dangerous characters
	sanitized = html.EscapeString(sanitized)

	// Remove null bytes and control characters
	sanitized = s.removeControlCharacters(sanitized)

	// Trim whitespace
	sanitized = strings.TrimSpace(sanitized)

	return sanitized
}

// SanitizeEmail sanitizes and validates email format
func (s *Sanitizer) SanitizeEmail(email string) string {
	if email == "" {
		return email
	}

	// Convert to lowercase
	email = strings.ToLower(email)

	// Remove dangerous characters except @ and .
	emailRegex := regexp.MustCompile(`[^a-zA-Z0-9@._-]`)
	email = emailRegex.ReplaceAllString(email, "")

	// Trim whitespace
	email = strings.TrimSpace(email)

	return email
}

// SanitizeName sanitizes name fields (allows letters, spaces, hyphens, apostrophes)
func (s *Sanitizer) SanitizeName(name string) string {
	if name == "" {
		return name
	}

	// Remove HTML tags and scripts
	sanitized := s.scriptTagRegex.ReplaceAllString(name, "")
	sanitized = s.htmlTagRegex.ReplaceAllString(sanitized, "")

	// Allow only letters, spaces, hyphens, and apostrophes
	nameRegex := regexp.MustCompile(`[^a-zA-Z\s\-']`)
	sanitized = nameRegex.ReplaceAllString(sanitized, "")

	// Remove multiple spaces
	spaceRegex := regexp.MustCompile(`\s+`)
	sanitized = spaceRegex.ReplaceAllString(sanitized, " ")

	// Trim whitespace
	sanitized = strings.TrimSpace(sanitized)

	return sanitized
}

// SanitizeAlphanumeric sanitizes input to allow only alphanumeric characters
func (s *Sanitizer) SanitizeAlphanumeric(input string) string {
	if input == "" {
		return input
	}

	// Allow only alphanumeric characters
	alphanumericRegex := regexp.MustCompile(`[^a-zA-Z0-9]`)
	sanitized := alphanumericRegex.ReplaceAllString(input, "")

	return sanitized
}

// DetectSQLInjection checks for potential SQL injection patterns
func (s *Sanitizer) DetectSQLInjection(input string) bool {
	if input == "" {
		return false
	}

	// Convert to lowercase for case-insensitive matching
	lowerInput := strings.ToLower(input)

	return s.sqlInjectionRegex.MatchString(lowerInput)
}

// DetectXSS checks for potential XSS patterns
func (s *Sanitizer) DetectXSS(input string) bool {
	if input == "" {
		return false
	}

	// Convert to lowercase for case-insensitive matching
	lowerInput := strings.ToLower(input)

	// Check for script tags
	if s.scriptTagRegex.MatchString(lowerInput) {
		return true
	}

	// Check for javascript: protocol
	if strings.Contains(lowerInput, "javascript:") {
		return true
	}

	// Check for event handlers
	eventHandlers := []string{"onload", "onerror", "onclick", "onmouseover", "onfocus", "onblur"}
	for _, handler := range eventHandlers {
		if strings.Contains(lowerInput, handler) {
			return true
		}
	}

	return false
}

// removeControlCharacters removes control characters and null bytes
func (s *Sanitizer) removeControlCharacters(input string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsControl(r) && r != '\n' && r != '\r' && r != '\t' {
			return -1 // Remove the character
		}
		return r
	}, input)
}

// SanitizeStruct sanitizes all string fields in a struct using reflection
func (s *Sanitizer) SanitizeStruct(data interface{}) interface{} {
	// This would require reflection to iterate through struct fields
	// For now, we'll implement specific sanitization in the middleware
	return data
}

// ValidateInputLength checks if input length is within acceptable limits
func (s *Sanitizer) ValidateInputLength(input string, maxLength int) bool {
	return len(input) <= maxLength
}

// Global sanitizer instance
var DefaultSanitizer = NewSanitizer()

// Convenience functions using the default sanitizer
func SanitizeString(input string) string {
	return DefaultSanitizer.SanitizeString(input)
}

func SanitizeEmail(email string) string {
	return DefaultSanitizer.SanitizeEmail(email)
}

func SanitizeName(name string) string {
	return DefaultSanitizer.SanitizeName(name)
}

func DetectSQLInjection(input string) bool {
	return DefaultSanitizer.DetectSQLInjection(input)
}

func DetectXSS(input string) bool {
	return DefaultSanitizer.DetectXSS(input)
}
