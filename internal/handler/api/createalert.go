package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *ApiHandler) CreateAlertHandler(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]string{"TODO": "WIP"})
}
