package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/venturarome/DaftWatch/internal/scraper"
)

func (h *ApiHandler) SearchHandler(c *gin.Context) {
	// TODO get prameters to compose Daft URL

	criteria := scraper.Criteria{
		SearchType: "rent",
		Location:   "dublin", //"dublin-9-dublin",
		Filters: []scraper.Filter{
			{
				Key:   "maxPrice",
				Value: "2000",
			},
			{
				Key:   "minBedrooms",
				Value: "3",
			},
		},
	}

	//resp := scraper.Scrape("/property-for-rent/dublin-9-dublin?rentalPrice_to=2000")
	resp := scraper.Scrape(criteria)
	c.JSON(http.StatusOK, resp)
}
