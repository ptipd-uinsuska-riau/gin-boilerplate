package health

import (
	"context"

	"gin-boilerplate/infrastructure/log"
	"gin-boilerplate/modules/primitive"

	"github.com/go-redis/redis"
)

type ServiceInterface interface {
	CheckUpTime(ctx context.Context) (resp primitive.HealthResponse, err error)
}

type Service struct {
	repository  RepositoryInterface
	redisClient *redis.Client
}

func NewService(repository RepositoryInterface, redisClient *redis.Client) ServiceInterface {
	return &Service{
		repository:  repository,
		redisClient: redisClient,
	}
}

func (s *Service) CheckUpTime(ctx context.Context) (primitive.HealthResponse, error) {
	// Check database connection
	errCheckDb := s.repository.CheckUpTimeDB(ctx)
	if errCheckDb != nil {
		log.Error(ctx, "Database health check failed:", errCheckDb)
		return primitive.HealthResponse{}, errCheckDb
	}
	databaseStatus := "healthy"

	// Check Redis connection (if enabled)
	var redisStatus string
	if s.redisClient != nil {
		errCheckRedis := s.redisClient.Ping().Err()
		if errCheckRedis != nil {
			log.Error(ctx, "Redis health check failed:", errCheckRedis)
			return primitive.HealthResponse{}, errCheckRedis
		}
		redisStatus = "healthy"
	} else {
		redisStatus = "disabled"
	}

	return primitive.HealthResponse{
		Database: databaseStatus,
		Redis:    redisStatus,
	}, nil
}
