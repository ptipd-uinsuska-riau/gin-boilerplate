// @title Gin Boilerplate API
// @version 1.0
// @description A clean and well-structured Go Gin boilerplate with authentication, user management, and essential features.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gin-boilerplate/boot"
	"gin-boilerplate/infrastructure/config"
	"gin-boilerplate/router"

	log "github.com/sirupsen/logrus"
)

func main() {
	// Parse command line flags
	flag.StringVar(&config.Env, "env", "local", "Environment configuration to use (local, dev, uat, prod)")
	flag.Parse()

	// Initialize application
	setup := boot.MakeHandler()
	defer func() {
		if setup.DB != nil {
			if err := setup.DB.Close(); err != nil {
				log.Errorf("Error closing database connection: %v", err)
			}
		}
	}()

	// Setup router
	handlerRouter := router.NewHandlerRouter(setup)
	app := handlerRouter.RouterWithMiddleware()

	// Configure server port
	port := fmt.Sprintf(":%d", config.Conf.Port)
	if config.Conf.Port == 0 {
		port = ":8080"
	}

	// Create HTTP server
	server := &http.Server{
		Addr:    port,
		Handler: app,
	}

	// Start server in a goroutine
	go func() {
		log.Infof("🚀 Server starting on port %s", port)
		log.Infof("📖 Swagger documentation available at: http://localhost%s/swagger/index.html", port)
		log.Infof("🌍 Environment: %s", config.Conf.Env)
		
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	// kill (no param) default sends syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("🛑 Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	select {
	case <-ctx.Done():
		log.Info("⏰ Timeout of 5 seconds reached")
	default:
		log.Info("✅ Server exited gracefully")
	}
}
