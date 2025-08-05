package middleware

import (
	"strings"

	"gin-boilerplate/infrastructure/errors"
	"gin-boilerplate/infrastructure/jwt"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware handles JWT authentication
type AuthMiddleware struct {
	jwtService *jwt.JWTService
}

// NewAuthMiddleware creates a new auth middleware
func NewAuthMiddleware(jwtService *jwt.JWTService) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService: jwtService,
	}
}

// RequireAuth middleware that requires authentication
func (am *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractToken(c)
		if token == "" {
			errors.AbortWithError(c, errors.NewUnauthorizedError("Authorization token required"))
			return
		}

		claims, err := am.jwtService.ValidateToken(token)
		if err != nil {
			errors.AbortWithError(c, errors.NewUnauthorizedError("Invalid or expired token"))
			return
		}

		// Set user info in context
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("claims", claims)

		c.Next()
	}
}

// OptionalAuth middleware that optionally checks authentication
func (am *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractToken(c)
		if token != "" {
			claims, err := am.jwtService.ValidateToken(token)
			if err == nil {
				// Set user info in context if token is valid
				c.Set("user_id", claims.UserID)
				c.Set("user_email", claims.Email)
				c.Set("claims", claims)
			}
		}

		c.Next()
	}
}

// extractToken extracts JWT token from Authorization header
func extractToken(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return ""
	}

	// Check for Bearer token
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	return parts[1]
}

// GetUserID gets user ID from context
func GetUserID(c *gin.Context) (string, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return "", false
	}
	return userID.(string), true
}

// GetUserEmail gets user email from context
func GetUserEmail(c *gin.Context) (string, bool) {
	email, exists := c.Get("user_email")
	if !exists {
		return "", false
	}
	return email.(string), true
}

// GetClaims gets JWT claims from context
func GetClaims(c *gin.Context) (*jwt.Claims, bool) {
	claims, exists := c.Get("claims")
	if !exists {
		return nil, false
	}
	return claims.(*jwt.Claims), true
}
