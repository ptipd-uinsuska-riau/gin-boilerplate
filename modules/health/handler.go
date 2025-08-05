package health

import (
	"net/http"

	"gin-boilerplate/infrastructure/httplib"
	"gin-boilerplate/infrastructure/log"

	"github.com/gin-gonic/gin"
)

type Http struct {
	serviceHealth ServiceInterface
}

func NewHttp(serviceHealth ServiceInterface) HttpInterface {
	return &Http{
		serviceHealth: serviceHealth,
	}
}

type HttpInterface interface {
	GroupHealth(group *gin.RouterGroup)
}

func (h *Http) GroupHealth(g *gin.RouterGroup) {
	g.GET("/ping", h.Ping)
	g.GET("/check", h.HealthCheckApi)
}

// Ping godoc
// @Summary Health check ping
// @Description Simple health check endpoint that returns pong
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} primitive.APIResponse "Successfully pinged"
// @Router /health/ping [get]
func (h *Http) Ping(c *gin.Context) {
	httplib.SetSuccessResponse(c, 1, http.StatusOK, "pong", nil)
}

// HealthCheckApi godoc
// @Summary Health check with system information
// @Description Comprehensive health check that returns system uptime and status
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} primitive.APIResponse "Successfully retrieved health status"
// @Failure 500 {object} primitive.APIErrorResponse "Internal server error"
// @Router /health/check [get]
func (h *Http) HealthCheckApi(ctx *gin.Context) {
	resp, err := h.serviceHealth.CheckUpTime(ctx)
	if err != nil {
		log.Error(ctx, "Health check failed:", err)
		httplib.SetErrorResponse(ctx, 0, http.StatusInternalServerError, err.Error())
		return
	}

	httplib.SetSuccessResponse(ctx, 1, http.StatusOK, "Health check successful", resp)
}
