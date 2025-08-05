package auth

import (
	"context"
	"time"

	"gin-boilerplate/infrastructure/errors"
	"gin-boilerplate/infrastructure/jwt"
	"gin-boilerplate/modules/primitive"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ServiceInterface interface {
	Register(ctx context.Context, req primitive.RegisterRequest) (*primitive.LoginResponse, error)
	Login(ctx context.Context, req primitive.LoginRequest) (*primitive.LoginResponse, error)
	GetProfile(ctx context.Context, userID uuid.UUID) (*primitive.UserInfo, error)
	UpdateProfile(ctx context.Context, userID uuid.UUID, req primitive.UpdateUserRequest) (*primitive.UserInfo, error)
	ChangePassword(ctx context.Context, userID uuid.UUID, req primitive.ChangePasswordRequest) error
}

type Service struct {
	repository RepositoryInterface
	jwtService *jwt.JWTService
}

func NewService(repository RepositoryInterface, jwtService *jwt.JWTService) ServiceInterface {
	return &Service{
		repository: repository,
		jwtService: jwtService,
	}
}

func (s *Service) Register(ctx context.Context, req primitive.RegisterRequest) (*primitive.LoginResponse, error) {
	// Check if email already exists
	exists, err := s.repository.EmailExists(ctx, req.Email)
	if err != nil {
		return nil, errors.NewInternalServerError("Failed to check email existence", err)
	}
	if exists {
		return nil, errors.NewConflictError(errors.ErrEmailExists)
	}

	// Create new user
	user := &User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	// Hash password
	if err := user.HashPassword(); err != nil {
		return nil, errors.NewInternalServerError("Failed to hash password", err)
	}

	// Save user
	if err := s.repository.CreateUser(ctx, user); err != nil {
		return nil, errors.NewInternalServerError("Failed to create user", err)
	}

	// Generate token
	token, err := s.jwtService.GenerateToken(user.ID.String(), user.Email)
	if err != nil {
		return nil, errors.NewInternalServerError("Failed to generate token", err)
	}

	// Calculate expiry time
	expiresAt := time.Now().Add(time.Hour * 24) // Default 24 hours

	return &primitive.LoginResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User:      user.ToUserInfo(),
	}, nil
}

func (s *Service) Login(ctx context.Context, req primitive.LoginRequest) (*primitive.LoginResponse, error) {
	// Get user by email
	user, err := s.repository.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewUnauthorizedError("Invalid email or password")
		}
		return nil, errors.NewInternalServerError("Failed to get user", err)
	}

	// Check password
	if !user.CheckPassword(req.Password) {
		return nil, errors.NewUnauthorizedError("Invalid email or password")
	}

	// Generate token
	token, err := s.jwtService.GenerateToken(user.ID.String(), user.Email)
	if err != nil {
		return nil, errors.NewInternalServerError("Failed to generate token", err)
	}

	// Calculate expiry time
	expiresAt := time.Now().Add(time.Hour * 24) // Default 24 hours

	return &primitive.LoginResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User:      user.ToUserInfo(),
	}, nil
}

func (s *Service) GetProfile(ctx context.Context, userID uuid.UUID) (*primitive.UserInfo, error) {
	user, err := s.repository.GetUserByID(ctx, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewNotFoundError(errors.ErrUserNotFound)
		}
		return nil, errors.NewInternalServerError("Failed to get user", err)
	}

	userInfo := user.ToUserInfo()
	return &userInfo, nil
}

func (s *Service) UpdateProfile(ctx context.Context, userID uuid.UUID, req primitive.UpdateUserRequest) (*primitive.UserInfo, error) {
	user, err := s.repository.GetUserByID(ctx, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewNotFoundError(errors.ErrUserNotFound)
		}
		return nil, errors.NewInternalServerError("Failed to get user", err)
	}

	// Update fields if provided
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		// Check if new email already exists
		if req.Email != user.Email {
			exists, err := s.repository.EmailExists(ctx, req.Email)
			if err != nil {
				return nil, errors.NewInternalServerError("Failed to check email existence", err)
			}
			if exists {
				return nil, errors.NewConflictError(errors.ErrEmailExists)
			}
		}
		user.Email = req.Email
	}

	// Save updated user
	if err := s.repository.UpdateUser(ctx, user); err != nil {
		return nil, errors.NewInternalServerError("Failed to update user", err)
	}

	userInfo := user.ToUserInfo()
	return &userInfo, nil
}

func (s *Service) ChangePassword(ctx context.Context, userID uuid.UUID, req primitive.ChangePasswordRequest) error {
	user, err := s.repository.GetUserByID(ctx, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.NewNotFoundError(errors.ErrUserNotFound)
		}
		return errors.NewInternalServerError("Failed to get user", err)
	}

	// Check current password
	if !user.CheckPassword(req.CurrentPassword) {
		return errors.NewUnauthorizedError(errors.ErrInvalidPassword)
	}

	// Update password
	user.Password = req.NewPassword
	if err := user.HashPassword(); err != nil {
		return errors.NewInternalServerError("Failed to hash password", err)
	}

	// Save updated user
	if err := s.repository.UpdateUser(ctx, user); err != nil {
		return errors.NewInternalServerError("Failed to update password", err)
	}

	return nil
}
