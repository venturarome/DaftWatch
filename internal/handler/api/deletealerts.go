package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *ApiHandler) DeleteAlertsHandler(c *gin.Context) {
	// TODO pass Property instance as parameter
	c.JSON(http.StatusOK, h.dbClient.DeleteAlerts())
}
