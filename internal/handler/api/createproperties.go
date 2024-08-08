package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *ApiHandler) CreatePropertiesHandler(c *gin.Context) {
	// TODO pass Property instances as parameters
	//c.JSON(http.StatusOK, h.dbClient.CreateProperties())
	c.JSON(http.StatusOK, map[string]string{"TODO": "WIP"})
}
