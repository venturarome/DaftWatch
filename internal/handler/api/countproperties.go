package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *ApiHandler) CountPropertiesHandler(c *gin.Context) {
	// TODO pass Property instances as parameters
	c.JSON(http.StatusOK, h.dbClient.CountProperties())
}
