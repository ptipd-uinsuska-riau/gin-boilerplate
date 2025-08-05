package redis

import (
	"fmt"
	"time"

	"gin-boilerplate/infrastructure/config"

	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
)

type LibInterface interface {
	Set(key string, value interface{}, expiration time.Duration) error
	Get(key string) (string, error)
	Del(key string) error
	Exists(key string) (bool, error)
	SetNX(key string, value interface{}, expiration time.Duration) (bool, error)
	Expire(key string, expiration time.Duration) error
}

type RedisLib struct {
	client *redis.Client
}

func NewRedisClient(config *config.Config) (*redis.Client, error) {
	if !config.Redis.EnableRedis {
		return nil, nil
	}

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Redis.Host, config.Redis.Port),
		Password: config.Redis.Password,
		DB:       config.Redis.DB,
	})

	// Test the connection
	_, err := client.Ping().Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	log.Info("Successfully connected to Redis")
	return client, nil
}

func NewRedisLibInterface(client *redis.Client) (LibInterface, error) {
	if client == nil {
		return nil, fmt.Errorf("redis client is nil")
	}

	return &RedisLib{
		client: client,
	}, nil
}

func (r *RedisLib) Set(key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(key, value, expiration).Err()
}

func (r *RedisLib) Get(key string) (string, error) {
	return r.client.Get(key).Result()
}

func (r *RedisLib) Del(key string) error {
	return r.client.Del(key).Err()
}

func (r *RedisLib) Exists(key string) (bool, error) {
	result, err := r.client.Exists(key).Result()
	return result > 0, err
}

func (r *RedisLib) SetNX(key string, value interface{}, expiration time.Duration) (bool, error) {
	return r.client.SetNX(key, value, expiration).Result()
}

func (r *RedisLib) Expire(key string, expiration time.Duration) error {
	return r.client.Expire(key, expiration).Err()
}
