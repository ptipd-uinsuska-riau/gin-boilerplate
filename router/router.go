package router

import (
	"net/http"
	"os"
	"time"

	"gin-boilerplate/boot"
	"gin-boilerplate/infrastructure/config"
	"gin-boilerplate/infrastructure/errors"
	"gin-boilerplate/infrastructure/httplib"
	"gin-boilerplate/infrastructure/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type HandlerRouter struct {
	Setup boot.HandlerSetup
}

type InterfaceRouter interface {
	RouterWithMiddleware() *gin.Engine
}

func NewHandlerRouter(setup boot.HandlerSetup) InterfaceRouter {
	return &HandlerRouter{
		Setup: setup,
	}
}

func notFoundHandler(c *gin.Context) {
	httplib.SetErrorResponse(c, http.StatusNotFound, http.StatusNotFound, "Route not found")
}

func methodNotAllowedHandler(c *gin.Context) {
	httplib.SetErrorResponse(c, http.StatusMethodNotAllowed, http.StatusMethodNotAllowed, "Method not allowed")
}

func (hr *HandlerRouter) RouterWithMiddleware() *gin.Engine {
	// Create new Gin instance
	r := gin.New()

	// Handle OPTIONS requests
	r.OPTIONS("*any", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// Configure CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Configure this properly for production
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Use recovery middleware
	r.Use(gin.Recovery())

	// Use centralized error handler
	r.Use(errors.ErrorHandler())

	// Security middlewares
	if config.Conf.Security.EnableSanitization {
		// Add input sanitization middleware
		sanitizerConfig := middleware.SanitizerConfig{
			EnableXSSDetection:          config.Conf.Security.EnableXSSDetection,
			EnableSQLInjectionDetection: config.Conf.Security.EnableSQLInjectionDetection,
			MaxStringLength:             config.Conf.Security.MaxStringLength,
			StrictMode:                  config.Conf.Security.StrictMode,
			SkipFields:                  []string{"password", "token", "secret"},
		}
		r.Use(middleware.SanitizerMiddleware(sanitizerConfig))

		// Add header and query parameter sanitization
		r.Use(middleware.HeaderSanitizerMiddleware())
		r.Use(middleware.QueryParamSanitizerMiddleware())
	}

	// Add request size limiting middleware
	requestLimiterConfig := middleware.RequestLimiterConfig{
		MaxBodySize:    config.Conf.Security.MaxBodySize,
		MaxFormFields:  config.Conf.Security.MaxFormFields,
		MaxFormMemory:  config.Conf.Security.MaxFormMemory,
		MaxFormFiles:   config.Conf.Security.MaxFormFiles,
		MaxHeaderSize:  config.Conf.Security.MaxHeaderSize,
		MaxURLLength:   config.Conf.Security.MaxURLLength,
		MaxQueryParams: config.Conf.Security.MaxQueryParams,
		SkipPaths:      []string{"/api/v1/health"},
	}
	r.Use(middleware.RequestLimiterMiddleware(requestLimiterConfig))

	// Configure logging
	if config.Conf.LogMode {
		// Log to file
		f, err := os.OpenFile("gin-boilerplate.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		r.Use(gin.LoggerWithConfig(gin.LoggerConfig{
			Output: f,
		}))
	} else {
		// Log to stdout
		r.Use(gin.LoggerWithConfig(gin.LoggerConfig{
			Output: gin.DefaultWriter,
		}))
	}

	// Set custom handlers for 404 and 405
	r.NoRoute(notFoundHandler)
	r.NoMethod(methodNotAllowedHandler)

	// Add rate limiting middleware
	r.Use(func(c *gin.Context) {
		clientID := getClientID(c)
		if !hr.Setup.RateLimiter.Allow(clientID) {
			errors.AbortWithError(c, errors.NewRateLimitError("Rate limit exceeded"))
			return
		}
		c.Next()
	})

	// API routes
	api := r.Group("/api")
	v1 := api.Group("/v1")

	// Health check routes
	healthGroup := v1.Group("/health")
	hr.Setup.HealthHttp.GroupHealth(healthGroup)

	// Authentication routes
	authGroup := v1.Group("/auth")
	hr.Setup.AuthHttp.GroupAuth(authGroup)

	// User management routes
	userGroup := v1.Group("/users")
	hr.Setup.UserHttp.GroupUser(userGroup)

	// File upload routes (example with file size limiting)
	uploadGroup := v1.Group("/upload")
	uploadGroup.Use(middleware.FileSizeLimiterMiddleware(config.Conf.Security.MaxFileSize))
	{
		uploadGroup.POST("/file", func(c *gin.Context) {
			httplib.SetSuccessResponse(c, 1, http.StatusOK, "File upload endpoint", map[string]string{
				"message": "File upload functionality can be implemented here",
				"maxSize": "Check security.maxFileSize in config",
			})
		})
	}

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Root endpoint
	r.GET("/", func(c *gin.Context) {
		httplib.SetSuccessResponse(c, 1, http.StatusOK, "Gin Boilerplate API", map[string]string{
			"version":     "1.0.0",
			"environment": config.Conf.Env,
			"swagger":     "/swagger/index.html",
		})
	})

	return r
}

// getClientID extracts client identifier for rate limiting
func getClientID(c *gin.Context) string {
	// Try to get user ID from context first (for authenticated users)
	if userID, exists := c.Get("user_id"); exists {
		return "user:" + userID.(string)
	}

	// Fall back to IP address
	return "ip:" + c.ClientIP()
}
