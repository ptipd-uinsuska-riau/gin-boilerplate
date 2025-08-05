package user

import (
	"context"

	"gin-boilerplate/infrastructure/errors"
	"gin-boilerplate/infrastructure/httplib"
	"gin-boilerplate/modules/auth"
	"gin-boilerplate/modules/primitive"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ServiceInterface interface {
	GetUsers(ctx context.Context, params primitive.QueryParams) (*httplib.PaginatedResponse, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*primitive.UserInfo, error)
	CreateUser(ctx context.Context, req primitive.RegisterRequest) (*primitive.UserInfo, error)
	UpdateUser(ctx context.Context, id uuid.UUID, req primitive.UpdateUserRequest) (*primitive.UserInfo, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

type Service struct {
	repository RepositoryInterface
}

func NewService(repository RepositoryInterface) ServiceInterface {
	return &Service{
		repository: repository,
	}
}

func (s *Service) GetUsers(ctx context.Context, params primitive.QueryParams) (*httplib.PaginatedResponse, error) {
	// Set defaults
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit <= 0 {
		params.Limit = 10
	}

	offset := (params.Page - 1) * params.Limit

	users, total, err := s.repository.GetUsers(ctx, offset, params.Limit, params.Search)
	if err != nil {
		return nil, errors.NewInternalServerError("Failed to get users", err)
	}

	// Convert to UserInfo
	var userInfos []primitive.UserInfo
	for _, user := range users {
		userInfos = append(userInfos, user.ToUserInfo())
	}

	pagination := httplib.CalculatePagination(params.Page, params.Limit, total)

	return &httplib.PaginatedResponse{
		Success:    true,
		Code:       1,
		Message:    "Users retrieved successfully",
		Data:       userInfos,
		Pagination: pagination,
	}, nil
}

func (s *Service) GetUserByID(ctx context.Context, id uuid.UUID) (*primitive.UserInfo, error) {
	user, err := s.repository.GetUserByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewNotFoundError(errors.ErrUserNotFound)
		}
		return nil, errors.NewInternalServerError("Failed to get user", err)
	}

	userInfo := user.ToUserInfo()
	return &userInfo, nil
}

func (s *Service) CreateUser(ctx context.Context, req primitive.RegisterRequest) (*primitive.UserInfo, error) {
	// Check if email already exists
	exists, err := s.repository.EmailExists(ctx, req.Email, nil)
	if err != nil {
		return nil, errors.NewInternalServerError("Failed to check email existence", err)
	}
	if exists {
		return nil, errors.NewConflictError(errors.ErrEmailExists)
	}

	// Create new user
	user := &auth.User{
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

	userInfo := user.ToUserInfo()
	return &userInfo, nil
}

func (s *Service) UpdateUser(ctx context.Context, id uuid.UUID, req primitive.UpdateUserRequest) (*primitive.UserInfo, error) {
	user, err := s.repository.GetUserByID(ctx, id)
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
			exists, err := s.repository.EmailExists(ctx, req.Email, &user.ID)
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

func (s *Service) DeleteUser(ctx context.Context, id uuid.UUID) error {
	// Check if user exists
	_, err := s.repository.GetUserByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.NewNotFoundError(errors.ErrUserNotFound)
		}
		return errors.NewInternalServerError("Failed to get user", err)
	}

	// Delete user
	if err := s.repository.DeleteUser(ctx, id); err != nil {
		return errors.NewInternalServerError("Failed to delete user", err)
	}

	return nil
}
