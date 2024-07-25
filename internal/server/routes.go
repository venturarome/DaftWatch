package server

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/venturarome/DaftWatch/internal/handler/api"
	"github.com/venturarome/DaftWatch/internal/scraper"
)

func (s *Server) RegisterRoutes() http.Handler {
	router := gin.Default()

	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// When the number of routes grow, I'll consider grouping:
	// https://gin-gonic.com/docs/examples/grouping-routes/
	router.GET("/", s.HelloWorldHandler)

	router.GET("/health", s.healthHandler)

	router.GET("/search", api.Search)

	router.Static("/images", "./images")

	return router
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

	resp := scraper.Scrape("/property-for-rent/dublin-9-dublin?rentalPrice_to=2000")

	//respJson, _ := json.Marshal(resp)

	// ret := make(map[string]string)
	// ret["message"] = "Hola mi amor"
	c.JSON(http.StatusOK, resp)
}
