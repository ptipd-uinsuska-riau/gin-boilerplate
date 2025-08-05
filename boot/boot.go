package boot

import (
	"os"

	"gin-boilerplate/infrastructure/config"
	"gin-boilerplate/infrastructure/database"
	"gin-boilerplate/infrastructure/jwt"
	"gin-boilerplate/infrastructure/log"
	"gin-boilerplate/infrastructure/middleware"
	"gin-boilerplate/infrastructure/redis"
	"gin-boilerplate/modules/auth"
	"gin-boilerplate/modules/health"
	"gin-boilerplate/modules/user"
	"gin-boilerplate/utils"

	redisClient "github.com/go-redis/redis"
	logrus "github.com/sirupsen/logrus"
)

type HandlerSetup struct {
	// Infrastructure
	RateLimiter    *middleware.RateLimiter
	AuthMiddleware *middleware.AuthMiddleware
	JWTService     *jwt.JWTService

	// Modules
	HealthHttp health.HttpInterface
	AuthHttp   auth.HttpInterface
	UserHttp   user.HttpInterface

	// Database
	DB *database.DatabaseClient

	// Redis (optional)
	RedisClient *redisClient.Client
}

func MakeHandler() HandlerSetup {
	// Initialize configuration
	config.Initialize()

	// Initialize logger
	log.Init(config.Conf.LogFormat, config.Conf.LogLevel)

	var err error

	// Initialize Redis client (optional)
	var redisClientInstance *redisClient.Client
	if config.Conf.Redis.EnableRedis {
		redisClientInstance, err = redis.NewRedisClient(&config.Conf)
		if err != nil {
			logrus.Fatalf("Failed to initialize Redis: %v", err)
			os.Exit(1)
		}
	}

	// Initialize database
	db, err := database.NewDatabaseClient(&config.Conf)
	if err != nil {
		logrus.Fatalf("Failed to initialize database: %v", err)
		os.Exit(1)
	}

	// Auto-migrate database tables
	if err := db.AutoMigrate(&auth.User{}); err != nil {
		logrus.Fatalf("Failed to migrate database: %v", err)
		os.Exit(1)
	}

	// Initialize JWT service
	jwtService := jwt.NewJWTService(&config.Conf)

	// Initialize rate limiter
	interval := utils.StringUnitToDuration(config.Conf.Interval)
	rateLimiter := middleware.NewRateLimiter(int(config.Conf.Rate), interval)

	// Initialize auth middleware
	authMiddleware := middleware.NewAuthMiddleware(jwtService)

	// Initialize repositories
	healthRepository := health.NewRepository(db.DbConn)
	authRepository := auth.NewRepository(db.DbConn)
	userRepository := user.NewRepository(db.DbConn)

	// Initialize services
	healthService := health.NewService(healthRepository, redisClientInstance)
	authService := auth.NewService(authRepository, jwtService)
	userService := user.NewService(userRepository)

	// Initialize HTTP handlers
	healthHttp := health.NewHttp(healthService)
	authHttp := auth.NewHttp(authService, authMiddleware)
	userHttp := user.NewHttp(userService, authMiddleware)

	logrus.Info("✅ Application setup completed successfully")

	return HandlerSetup{
		// Infrastructure
		RateLimiter:    rateLimiter,
		AuthMiddleware: authMiddleware,
		JWTService:     jwtService,

		// Modules
		HealthHttp: healthHttp,
		AuthHttp:   authHttp,
		UserHttp:   userHttp,

		// Database
		DB: db,

		// Redis
		RedisClient: redisClientInstance,
	}
}
