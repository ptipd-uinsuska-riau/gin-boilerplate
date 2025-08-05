package database

import (
	"fmt"
	"time"

	"gin-boilerplate/infrastructure/config"

	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DatabaseClient struct {
	DbConn *gorm.DB
}

func NewDatabaseClient(config *config.Config) (*DatabaseClient, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		config.Postgres.Host,
		config.Postgres.User,
		config.Postgres.Password,
		config.Postgres.DBName,
		config.Postgres.Port,
		config.Postgres.SSLMode,
		config.Postgres.TimeZone,
	)

	// Configure GORM logger
	var gormLogger logger.Interface
	if config.LogMode {
		gormLogger = logger.Default.LogMode(logger.Info)
	} else {
		gormLogger = logger.Default.LogMode(logger.Silent)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying sql.DB to configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Configure connection pool
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Test the connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Info("Successfully connected to PostgreSQL database")

	return &DatabaseClient{
		DbConn: db,
	}, nil
}

func (dc *DatabaseClient) Close() error {
	sqlDB, err := dc.DbConn.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// AutoMigrate runs database migrations for the given models
func (dc *DatabaseClient) AutoMigrate(models ...interface{}) error {
	return dc.DbConn.AutoMigrate(models...)
}
