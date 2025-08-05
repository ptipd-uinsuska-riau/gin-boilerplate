package errors

import (
	"context"
	"net/http"

	"gin-boilerplate/infrastructure/httplib"
	logger "gin-boilerplate/infrastructure/log"

	"github.com/gin-gonic/gin"
)

// ErrorHandler is a centralized error handler middleware
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if there are any errors
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			// Handle different types of errors
			switch e := err.(type) {
			case *AppError:
				handleAppError(c, e)
			default:
				handleGenericError(c, err)
			}
		}
	}
}

func handleAppError(c *gin.Context, err *AppError) {
	// Log the error
	if err.Code >= 500 {
		logger.Error(context.Background(), "Application error:", err.Error())
	} else {
		logger.Warn(context.Background(), "Client error:", err.Error())
	}

	// Return error response
	httplib.SetErrorResponse(c, 0, err.Code, err.Message)
}

func handleGenericError(c *gin.Context, err error) {
	// Log the error
	logger.Error(context.Background(), "Unhandled error:", err.Error())

	// Return generic error response
	httplib.SetErrorResponse(c, 0, http.StatusInternalServerError, "Internal server error")
}

// AbortWithError aborts the request with an error
func AbortWithError(c *gin.Context, err *AppError) {
	c.Error(err)
	c.Abort()
}
