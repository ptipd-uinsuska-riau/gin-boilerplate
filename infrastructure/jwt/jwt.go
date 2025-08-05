package jwt

import (
	"fmt"
	"time"

	"gin-boilerplate/infrastructure/config"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Claims represents JWT claims
type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// JWTService handles JWT operations
type JWTService struct {
	secretKey string
	expiry    time.Duration
}

// NewJWTService creates a new JWT service
func NewJWTService(config *config.Config) *JWTService {
	expiry := time.Hour * 24 // default 24 hours
	if config.AccessTokenExpiry > 0 {
		switch config.AccessTokenExpiryUnit {
		case "minute":
			expiry = time.Minute * time.Duration(config.AccessTokenExpiry)
		case "hour":
			expiry = time.Hour * time.Duration(config.AccessTokenExpiry)
		case "day":
			expiry = time.Hour * 24 * time.Duration(config.AccessTokenExpiry)
		}
	}

	return &JWTService{
		secretKey: config.SecretKey,
		expiry:    expiry,
	}
}

// GenerateToken generates a new JWT token
func (j *JWTService) GenerateToken(userID, email string) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.New().String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(j.expiry)),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "gin-boilerplate",
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

// ValidateToken validates a JWT token and returns claims
func (j *JWTService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// RefreshToken generates a new token with extended expiry
func (j *JWTService) RefreshToken(tokenString string) (string, error) {
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}

	// Generate new token with same user info but new expiry
	return j.GenerateToken(claims.UserID, claims.Email)
}

// ExtractUserID extracts user ID from token
func (j *JWTService) ExtractUserID(tokenString string) (string, error) {
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}
	return claims.UserID, nil
}

// ExtractEmail extracts email from token
func (j *JWTService) ExtractEmail(tokenString string) (string, error) {
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}
	return claims.Email, nil
}
