package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/venturarome/DaftWatch/internal/scraper"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.GET("/", s.HelloWorldHandler)

	r.GET("/health", s.healthHandler)

	r.GET("/search", s.searchHandler)

	return r
}

func (s *Server) HelloWorldHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	c.JSON(http.StatusOK, resp)
}

func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.db.Health())
}

func (s *Server) searchHandler(c *gin.Context) {
	sc := &scraper.PropertyListScraper{}
	sc.Scrape("/property-for-rent/dublin-9-dublin?rentalPrice_to=2000")

	resp := make(map[string]string)
	resp["message"] = "Hola mi amor"

	c.JSON(http.StatusOK, resp)
}
