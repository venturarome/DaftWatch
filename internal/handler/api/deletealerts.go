package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *ApiHandler) DeleteAlertsHandler(c *gin.Context) {
	// For testing purposes
	// TODO pass Property instance as parameter
	c.JSON(http.StatusOK, map[string]string{"TODO": "WIP"})
}
