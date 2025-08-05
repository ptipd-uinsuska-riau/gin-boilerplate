package user

import (
	"net/http"

	"gin-boilerplate/infrastructure/errors"
	"gin-boilerplate/infrastructure/httplib"
	"gin-boilerplate/infrastructure/log"
	"gin-boilerplate/infrastructure/middleware"
	"gin-boilerplate/infrastructure/validator"
	"gin-boilerplate/modules/primitive"

	"github.com/gin-gonic/gin"
)

type Http struct {
	serviceUser    ServiceInterface
	authMiddleware *middleware.AuthMiddleware
}

func NewHttp(serviceUser ServiceInterface, authMiddleware *middleware.AuthMiddleware) HttpInterface {
	return &Http{
		serviceUser:    serviceUser,
		authMiddleware: authMiddleware,
	}
}

type HttpInterface interface {
	GroupUser(group *gin.RouterGroup)
}

func (h *Http) GroupUser(g *gin.RouterGroup) {
	// All user management routes require authentication
	g.Use(h.authMiddleware.RequireAuth())

	g.GET("", h.GetUsers)
	g.GET("/:id", h.GetUserByID)
	g.POST("", h.CreateUser)
	g.PUT("/:id", h.UpdateUser)
	g.DELETE("/:id", h.DeleteUser)
}

// GetUsers godoc
// @Summary Get all users
// @Description Get a paginated list of users with optional search
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param search query string false "Search by name or email"
// @Success 200 {object} httplib.PaginatedResponse "Users retrieved successfully"
// @Failure 401 {object} primitive.APIErrorResponse "Unauthorized"
// @Failure 500 {object} primitive.APIErrorResponse "Internal server error"
// @Router /users [get]
func (h *Http) GetUsers(c *gin.Context) {
	var params primitive.QueryParams
	if err := validator.BindQueryAndValidate(c, &params); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			errors.AbortWithError(c, appErr)
		} else {
			errors.AbortWithError(c, errors.NewBadRequestError(err.Error()))
		}
		return
	}

	resp, err := h.serviceUser.GetUsers(c, params)
	if err != nil {
		log.Error(c, "Get users failed:", err)
		if appErr, ok := err.(*errors.AppError); ok {
			errors.AbortWithError(c, appErr)
		} else {
			errors.AbortWithError(c, errors.NewInternalServerError("Get users failed", err))
		}
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetUserByID godoc
// @Summary Get user by ID
// @Description Get a specific user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID" format(uuid)
// @Success 200 {object} primitive.APIResponse{data=primitive.UserInfo} "User retrieved successfully"
// @Failure 400 {object} primitive.APIErrorResponse "Bad request"
// @Failure 401 {object} primitive.APIErrorResponse "Unauthorized"
// @Failure 404 {object} primitive.APIErrorResponse "User not found"
// @Failure 500 {object} primitive.APIErrorResponse "Internal server error"
// @Router /users/{id} [get]
func (h *Http) GetUserByID(c *gin.Context) {
	var params primitive.IDParam
	if err := validator.BindURIAndValidate(c, &params); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			errors.AbortWithError(c, appErr)
		} else {
			errors.AbortWithError(c, errors.NewBadRequestError(err.Error()))
		}
		return
	}

	resp, err := h.serviceUser.GetUserByID(c, params.ID)
	if err != nil {
		log.Error(c, "Get user by ID failed:", err)
		if appErr, ok := err.(*errors.AppError); ok {
			errors.AbortWithError(c, appErr)
		} else {
			errors.AbortWithError(c, errors.NewInternalServerError("Get user failed", err))
		}
		return
	}

	httplib.SetSuccessResponse(c, 1, http.StatusOK, "User retrieved successfully", resp)
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user (admin function)
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body primitive.RegisterRequest true "Create user request"
// @Success 201 {object} primitive.APIResponse{data=primitive.UserInfo} "User created successfully"
// @Failure 400 {object} primitive.APIErrorResponse "Bad request"
// @Failure 401 {object} primitive.APIErrorResponse "Unauthorized"
// @Failure 409 {object} primitive.APIErrorResponse "Email already exists"
// @Failure 422 {object} primitive.APIErrorResponse "Validation error"
// @Failure 500 {object} primitive.APIErrorResponse "Internal server error"
// @Router /users [post]
func (h *Http) CreateUser(c *gin.Context) {
	var req primitive.RegisterRequest
	if err := validator.BindAndValidate(c, &req); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			errors.AbortWithError(c, appErr)
		} else {
			errors.AbortWithError(c, errors.NewBadRequestError(err.Error()))
		}
		return
	}

	resp, err := h.serviceUser.CreateUser(c, req)
	if err != nil {
		log.Error(c, "Create user failed:", err)
		if appErr, ok := err.(*errors.AppError); ok {
			errors.AbortWithError(c, appErr)
		} else {
			errors.AbortWithError(c, errors.NewInternalServerError("Create user failed", err))
		}
		return
	}

	httplib.SetSuccessResponse(c, 1, http.StatusCreated, "User created successfully", resp)
}

// UpdateUser godoc
// @Summary Update user
// @Description Update a user's information
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID" format(uuid)
// @Param request body primitive.UpdateUserRequest true "Update user request"
// @Success 200 {object} primitive.APIResponse{data=primitive.UserInfo} "User updated successfully"
// @Failure 400 {object} primitive.APIErrorResponse "Bad request"
// @Failure 401 {object} primitive.APIErrorResponse "Unauthorized"
// @Failure 404 {object} primitive.APIErrorResponse "User not found"
// @Failure 409 {object} primitive.APIErrorResponse "Email already exists"
// @Failure 422 {object} primitive.APIErrorResponse "Validation error"
// @Failure 500 {object} primitive.APIErrorResponse "Internal server error"
// @Router /users/{id} [put]
func (h *Http) UpdateUser(c *gin.Context) {
	var params primitive.IDParam
	if err := validator.BindURIAndValidate(c, &params); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			errors.AbortWithError(c, appErr)
		} else {
			errors.AbortWithError(c, errors.NewBadRequestError(err.Error()))
		}
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

	resp, err := h.serviceUser.UpdateUser(c, params.ID, req)
	if err != nil {
		log.Error(c, "Update user failed:", err)
		if appErr, ok := err.(*errors.AppError); ok {
			errors.AbortWithError(c, appErr)
		} else {
			errors.AbortWithError(c, errors.NewInternalServerError("Update user failed", err))
		}
		return
	}

	httplib.SetSuccessResponse(c, 1, http.StatusOK, "User updated successfully", resp)
}

// DeleteUser godoc
// @Summary Delete user
// @Description Delete a user by ID
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID" format(uuid)
// @Success 200 {object} primitive.APIResponse "User deleted successfully"
// @Failure 400 {object} primitive.APIErrorResponse "Bad request"
// @Failure 401 {object} primitive.APIErrorResponse "Unauthorized"
// @Failure 404 {object} primitive.APIErrorResponse "User not found"
// @Failure 500 {object} primitive.APIErrorResponse "Internal server error"
// @Router /users/{id} [delete]
func (h *Http) DeleteUser(c *gin.Context) {
	var params primitive.IDParam
	if err := validator.BindURIAndValidate(c, &params); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			errors.AbortWithError(c, appErr)
		} else {
			errors.AbortWithError(c, errors.NewBadRequestError(err.Error()))
		}
		return
	}

	err := h.serviceUser.DeleteUser(c, params.ID)
	if err != nil {
		log.Error(c, "Delete user failed:", err)
		if appErr, ok := err.(*errors.AppError); ok {
			errors.AbortWithError(c, appErr)
		} else {
			errors.AbortWithError(c, errors.NewInternalServerError("Delete user failed", err))
		}
		return
	}

	httplib.SetSuccessResponse(c, 1, http.StatusOK, "User deleted successfully", nil)
}
