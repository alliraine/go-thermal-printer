package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthController struct {
}

func NewHealthController(group *gin.RouterGroup) {
	controller := &HealthController{}

	group.GET("/health", controller.getHealthStatusHandler)
}

func (hc *HealthController) getHealthStatusHandler(c *gin.Context) {
	c.Status(http.StatusNoContent)
}
