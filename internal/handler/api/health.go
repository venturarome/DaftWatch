package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *ApiHandler) HealthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, h.dbClient.Health())
}
