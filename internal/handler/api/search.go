package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/venturarome/DaftWatch/internal/scraper"
)

func Search(c *gin.Context) {
	resp := scraper.Scrape("/property-for-rent/dublin-9-dublin?rentalPrice_to=2000")
	c.JSON(http.StatusOK, resp)
}
