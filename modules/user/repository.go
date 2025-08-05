package user

import (
	"context"

	"gin-boilerplate/modules/auth"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RepositoryInterface interface {
	GetUsers(ctx context.Context, offset, limit int, search string) ([]*auth.User, int64, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*auth.User, error)
	CreateUser(ctx context.Context, user *auth.User) error
	UpdateUser(ctx context.Context, user *auth.User) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
	EmailExists(ctx context.Context, email string, excludeID *uuid.UUID) (bool, error)
}

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) RepositoryInterface {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetUsers(ctx context.Context, offset, limit int, search string) ([]*auth.User, int64, error) {
	var users []*auth.User
	var total int64

	query := r.db.WithContext(ctx).Model(&auth.User{})

	// Apply search filter
	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("name ILIKE ? OR email ILIKE ?", searchPattern, searchPattern)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get users with pagination
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *Repository) GetUserByID(ctx context.Context, id uuid.UUID) (*auth.User, error) {
	var user auth.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) CreateUser(ctx context.Context, user *auth.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *Repository) UpdateUser(ctx context.Context, user *auth.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *Repository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&auth.User{}, id).Error
}

func (r *Repository) EmailExists(ctx context.Context, email string, excludeID *uuid.UUID) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&auth.User{}).Where("email = ?", email)
	
	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}
	
	err := query.Count(&count).Error
	return count > 0, err
}
