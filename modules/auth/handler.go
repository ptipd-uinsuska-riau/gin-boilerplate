package auth

import (
	"net/http"

	"gin-boilerplate/infrastructure/errors"
	"gin-boilerplate/infrastructure/httplib"
	"gin-boilerplate/infrastructure/log"
	"gin-boilerplate/infrastructure/middleware"
	"gin-boilerplate/infrastructure/validator"
	"gin-boilerplate/modules/primitive"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Http struct {
	serviceAuth    ServiceInterface
	authMiddleware *middleware.AuthMiddleware
}

func NewHttp(serviceAuth ServiceInterface, authMiddleware *middleware.AuthMiddleware) HttpInterface {
	return &Http{
		serviceAuth:    serviceAuth,
		authMiddleware: authMiddleware,
	}
}

type HttpInterface interface {
	GroupAuth(group *gin.RouterGroup)
}

func (h *Http) GroupAuth(g *gin.RouterGroup) {
	g.POST("/register", h.Register)
	g.POST("/login", h.Login)
	
	// Protected routes
	protected := g.Group("")
	protected.Use(h.authMiddleware.RequireAuth())
	{
		protected.GET("/profile", h.GetProfile)
		protected.PUT("/profile", h.UpdateProfile)
		protected.POST("/change-password", h.ChangePassword)
	}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with name, email, and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body primitive.RegisterRequest true "Registration request"
// @Success 201 {object} primitive.APIResponse{data=primitive.LoginResponse} "User registered successfully"
// @Failure 400 {object} primitive.APIErrorResponse "Bad request"
// @Failure 409 {object} primitive.APIErrorResponse "Email already exists"
// @Failure 422 {object} primitive.APIErrorResponse "Validation error"
// @Failure 500 {object} primitive.APIErrorResponse "Internal server error"
// @Router /auth/register [post]
func (h *Http) Register(c *gin.Context) {
	var req primitive.RegisterRequest
	
	if err := validator.BindAndValidate(c, &req); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			errors.AbortWithError(c, appErr)
		} else {
			errors.AbortWithError(c, errors.NewBadRequestError(err.Error()))
		}
		return
	}

	resp, err := h.serviceAuth.Register(c, req)
	if err != nil {
		log.Error(c, "Registration failed:", err)
		if appErr, ok := err.(*errors.AppError); ok {
			errors.AbortWithError(c, appErr)
		} else {
			errors.AbortWithError(c, errors.NewInternalServerError("Registration failed", err))
		}
		return
	}

	httplib.SetSuccessResponse(c, 1, http.StatusCreated, "User registered successfully", resp)
}

// Login godoc
// @Summary Login user
// @Description Authenticate user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body primitive.LoginRequest true "Login request"
// @Success 200 {object} primitive.APIResponse{data=primitive.LoginResponse} "Login successful"
// @Failure 400 {object} primitive.APIErrorResponse "Bad request"
// @Failure 401 {object} primitive.APIErrorResponse "Invalid credentials"
// @Failure 422 {object} primitive.APIErrorResponse "Validation error"
// @Failure 500 {object} primitive.APIErrorResponse "Internal server error"
// @Router /auth/login [post]
func (h *Http) Login(c *gin.Context) {
	var req primitive.LoginRequest
	
	if err := validator.BindAndValidate(c, &req); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			errors.AbortWithError(c, appErr)
		} else {
			errors.AbortWithError(c, errors.NewBadRequestError(err.Error()))
		}
		return
	}

	resp, err := h.serviceAuth.Login(c, req)
	if err != nil {
		log.Error(c, "Login failed:", err)
		if appErr, ok := err.(*errors.AppError); ok {
			errors.AbortWithError(c, appErr)
		} else {
			errors.AbortWithError(c, errors.NewInternalServerError("Login failed", err))
		}
		return
	}

	httplib.SetSuccessResponse(c, 1, http.StatusOK, "Login successful", resp)
}

// GetProfile godoc
// @Summary Get user profile
// @Description Get current user's profile information
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} primitive.APIResponse{data=primitive.UserInfo} "Profile retrieved successfully"
// @Failure 401 {object} primitive.APIErrorResponse "Unauthorized"
// @Failure 404 {object} primitive.APIErrorResponse "User not found"
// @Failure 500 {object} primitive.APIErrorResponse "Internal server error"
// @Router /auth/profile [get]
func (h *Http) GetProfile(c *gin.Context) {
	userIDStr, exists := middleware.GetUserID(c)
	if !exists {
		errors.AbortWithError(c, errors.NewUnauthorizedError("User ID not found in context"))
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		errors.AbortWithError(c, errors.NewBadRequestError("Invalid user ID"))
		return
	}

	resp, err := h.serviceAuth.GetProfile(c, userID)
	if err != nil {
		log.Error(c, "Get profile failed:", err)
		if appErr, ok := err.(*errors.AppError); ok {
			errors.AbortWithError(c, appErr)
		} else {
			errors.AbortWithError(c, errors.NewInternalServerError("Get profile failed", err))
		}
		return
	}

	httplib.SetSuccessResponse(c, 1, http.StatusOK, "Profile retrieved successfully", resp)
}

// UpdateProfile godoc
// @Summary Update user profile
// @Description Update current user's profile information
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body primitive.UpdateUserRequest true "Update profile request"
// @Success 200 {object} primitive.APIResponse{data=primitive.UserInfo} "Profile updated successfully"
// @Failure 400 {object} primitive.APIErrorResponse "Bad request"
// @Failure 401 {object} primitive.APIErrorResponse "Unauthorized"
// @Failure 404 {object} primitive.APIErrorResponse "User not found"
// @Failure 409 {object} primitive.APIErrorResponse "Email already exists"
// @Failure 422 {object} primitive.APIErrorResponse "Validation error"
// @Failure 500 {object} primitive.APIErrorResponse "Internal server error"
// @Router /auth/profile [put]
func (h *Http) UpdateProfile(c *gin.Context) {
	userIDStr, exists := middleware.GetUserID(c)
	if !exists {
		errors.AbortWithError(c, errors.NewUnauthorizedError("User ID not found in context"))
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		errors.AbortWithError(c, errors.NewBadRequestError("Invalid user ID"))
		return
	}

	var req primitive.UpdateUserRequest
	if err := validator.BindAndValidate(c, &req); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			errors.AbortWithError(c, appErr)
		} else {
			errors.AbortWithError(c, errors.NewBadRequestError(err.Error()))
		}
		return
	}

	resp, err := h.serviceAuth.UpdateProfile(c, userID, req)
	if err != nil {
		log.Error(c, "Update profile failed:", err)
		if appErr, ok := err.(*errors.AppError); ok {
			errors.AbortWithError(c, appErr)
		} else {
			errors.AbortWithError(c, errors.NewInternalServerError("Update profile failed", err))
		}
		return
	}

	httplib.SetSuccessResponse(c, 1, http.StatusOK, "Profile updated successfully", resp)
}

// ChangePassword godoc
// @Summary Change user password
// @Description Change current user's password
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body primitive.ChangePasswordRequest true "Change password request"
// @Success 200 {object} primitive.APIResponse "Password changed successfully"
// @Failure 400 {object} primitive.APIErrorResponse "Bad request"
// @Failure 401 {object} primitive.APIErrorResponse "Unauthorized or invalid current password"
// @Failure 404 {object} primitive.APIErrorResponse "User not found"
// @Failure 422 {object} primitive.APIErrorResponse "Validation error"
// @Failure 500 {object} primitive.APIErrorResponse "Internal server error"
// @Router /auth/change-password [post]
func (h *Http) ChangePassword(c *gin.Context) {
	userIDStr, exists := middleware.GetUserID(c)
	if !exists {
		errors.AbortWithError(c, errors.NewUnauthorizedError("User ID not found in context"))
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		errors.AbortWithError(c, errors.NewBadRequestError("Invalid user ID"))
		return
	}

	var req primitive.ChangePasswordRequest
	if err := validator.BindAndValidate(c, &req); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			errors.AbortWithError(c, appErr)
		} else {
			errors.AbortWithError(c, errors.NewBadRequestError(err.Error()))
		}
		return
	}

	err = h.serviceAuth.ChangePassword(c, userID, req)
	if err != nil {
		log.Error(c, "Change password failed:", err)
		if appErr, ok := err.(*errors.AppError); ok {
			errors.AbortWithError(c, appErr)
		} else {
			errors.AbortWithError(c, errors.NewInternalServerError("Change password failed", err))
		}
		return
	}

	httplib.SetSuccessResponse(c, 1, http.StatusOK, "Password changed successfully", nil)
}
