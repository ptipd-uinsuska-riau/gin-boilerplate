package sanitizer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSanitizeString(t *testing.T) {
	sanitizer := NewSanitizer()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Normal string should remain unchanged",
			input:    "Hello World",
			expected: "Hello World",
		},
		{
			name:     "HTML tags should be removed",
			input:    "Hello <b>World</b>",
			expected: "Hello World",
		},
		{
			name:     "Script tags should be removed",
			input:    "Hello <script>alert('xss')</script> World",
			expected: "Hello  World",
		},
		{
			name:     "Dangerous characters should be escaped",
			input:    "Hello & <World>",
			expected: "Hello &amp;",
		},
		{
			name:     "Empty string should remain empty",
			input:    "",
			expected: "",
		},
		{
			name:     "Whitespace should be trimmed",
			input:    "  Hello World  ",
			expected: "Hello World",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sanitizer.SanitizeString(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSanitizeEmail(t *testing.T) {
	sanitizer := NewSanitizer()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Valid email should be lowercased",
			input:    "John.Doe@Example.COM",
			expected: "john.doe@example.com",
		},
		{
			name:     "Email with dangerous characters should be cleaned",
			input:    "john<script>@example.com",
			expected: "johnscript@example.com",
		},
		{
			name:     "Empty email should remain empty",
			input:    "",
			expected: "",
		},
		{
			name:     "Email with spaces should be trimmed",
			input:    "  john@example.com  ",
			expected: "john@example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sanitizer.SanitizeEmail(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSanitizeName(t *testing.T) {
	sanitizer := NewSanitizer()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Normal name should remain unchanged",
			input:    "John Doe",
			expected: "John Doe",
		},
		{
			name:     "Name with apostrophe should be preserved",
			input:    "O'Connor",
			expected: "O'Connor",
		},
		{
			name:     "Name with hyphen should be preserved",
			input:    "Mary-Jane",
			expected: "Mary-Jane",
		},
		{
			name:     "Name with numbers should be cleaned",
			input:    "John123 Doe",
			expected: "John Doe",
		},
		{
			name:     "Name with HTML should be cleaned",
			input:    "John <script>alert('xss')</script> Doe",
			expected: "John Doe",
		},
		{
			name:     "Multiple spaces should be normalized",
			input:    "John    Doe",
			expected: "John Doe",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sanitizer.SanitizeName(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSanitizeAlphanumeric(t *testing.T) {
	sanitizer := NewSanitizer()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Alphanumeric string should remain unchanged",
			input:    "abc123",
			expected: "abc123",
		},
		{
			name:     "String with special characters should be cleaned",
			input:    "abc!@#123",
			expected: "abc123",
		},
		{
			name:     "Empty string should remain empty",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sanitizer.SanitizeAlphanumeric(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDetectSQLInjection(t *testing.T) {
	sanitizer := NewSanitizer()

	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "Normal string should not be detected",
			input:    "Hello World",
			expected: false,
		},
		{
			name:     "SQL injection with UNION should be detected",
			input:    "1' UNION SELECT * FROM users --",
			expected: true,
		},
		{
			name:     "SQL injection with DROP should be detected",
			input:    "'; DROP TABLE users; --",
			expected: true,
		},
		{
			name:     "SQL injection with SELECT should be detected",
			input:    "1' OR '1'='1",
			expected: true,
		},
		{
			name:     "Case insensitive detection",
			input:    "1' or '1'='1",
			expected: true,
		},
		{
			name:     "Empty string should not be detected",
			input:    "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sanitizer.DetectSQLInjection(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDetectXSS(t *testing.T) {
	sanitizer := NewSanitizer()

	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "Normal string should not be detected",
			input:    "Hello World",
			expected: false,
		},
		{
			name:     "Script tag should be detected",
			input:    "<script>alert('xss')</script>",
			expected: true,
		},
		{
			name:     "JavaScript protocol should be detected",
			input:    "javascript:alert('xss')",
			expected: true,
		},
		{
			name:     "Event handler should be detected",
			input:    "<img onload='alert(1)'>",
			expected: true,
		},
		{
			name:     "Case insensitive detection",
			input:    "<SCRIPT>alert('xss')</SCRIPT>",
			expected: true,
		},
		{
			name:     "Empty string should not be detected",
			input:    "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sanitizer.DetectXSS(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestValidateInputLength(t *testing.T) {
	sanitizer := NewSanitizer()

	tests := []struct {
		name      string
		input     string
		maxLength int
		expected  bool
	}{
		{
			name:      "String within limit should be valid",
			input:     "Hello",
			maxLength: 10,
			expected:  true,
		},
		{
			name:      "String exceeding limit should be invalid",
			input:     "Hello World",
			maxLength: 5,
			expected:  false,
		},
		{
			name:      "Empty string should be valid",
			input:     "",
			maxLength: 10,
			expected:  true,
		},
		{
			name:      "String equal to limit should be valid",
			input:     "Hello",
			maxLength: 5,
			expected:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sanitizer.ValidateInputLength(tt.input, tt.maxLength)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestConvenienceFunctions(t *testing.T) {
	// Test global convenience functions
	result := SanitizeString("<script>alert('test')</script>Hello")
	assert.NotContains(t, result, "<script>")
	assert.Contains(t, result, "Hello")

	email := SanitizeEmail("TEST@EXAMPLE.COM")
	assert.Equal(t, "test@example.com", email)

	name := SanitizeName("John123 Doe")
	assert.Equal(t, "John Doe", name)

	assert.True(t, DetectSQLInjection("'; DROP TABLE users; --"))
	assert.True(t, DetectXSS("<script>alert('xss')</script>"))
}
