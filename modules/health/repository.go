package health

import (
	"context"

	"gorm.io/gorm"
)

type RepositoryInterface interface {
	CheckUpTimeDB(ctx context.Context) error
}

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) RepositoryInterface {
	return &Repository{
		db: db,
	}
}

func (r *Repository) CheckUpTimeDB(ctx context.Context) error {
	sqlDB, err := r.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.PingContext(ctx)
}
