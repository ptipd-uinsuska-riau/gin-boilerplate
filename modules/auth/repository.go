package auth

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RepositoryInterface interface {
	CreateUser(ctx context.Context, user *User) error
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*User, error)
	UpdateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
	EmailExists(ctx context.Context, email string) (bool, error)
}

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) RepositoryInterface {
	return &Repository{
		db: db,
	}
}

func (r *Repository) CreateUser(ctx context.Context, user *User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) GetUserByID(ctx context.Context, id uuid.UUID) (*User, error) {
	var user User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) UpdateUser(ctx context.Context, user *User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *Repository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&User{}, id).Error
}

func (r *Repository) EmailExists(ctx context.Context, email string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}
