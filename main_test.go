package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"gin-boilerplate/boot"
	"gin-boilerplate/infrastructure/config"
	"gin-boilerplate/router"

	"github.com/stretchr/testify/assert"
)

func TestHealthEndpoint(t *testing.T) {
	// Set test environment
	config.Env = "test"
	
	// Initialize test setup
	setup := boot.MakeHandler()
	handlerRouter := router.NewHandlerRouter(setup)
	app := handlerRouter.RouterWithMiddleware()

	// Create test request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/health/ping", nil)
	
	// Perform request
	app.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "pong")
}

func TestRootEndpoint(t *testing.T) {
	// Set test environment
	config.Env = "test"
	
	// Initialize test setup
	setup := boot.MakeHandler()
	handlerRouter := router.NewHandlerRouter(setup)
	app := handlerRouter.RouterWithMiddleware()

	// Create test request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	
	// Perform request
	app.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Gin Boilerplate API")
}

func TestNotFoundEndpoint(t *testing.T) {
	// Set test environment
	config.Env = "test"
	
	// Initialize test setup
	setup := boot.MakeHandler()
	handlerRouter := router.NewHandlerRouter(setup)
	app := handlerRouter.RouterWithMiddleware()

	// Create test request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/nonexistent", nil)
	
	// Perform request
	app.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "Route not found")
}
