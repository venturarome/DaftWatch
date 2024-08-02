package api

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/venturarome/DaftWatch/internal/database"
)

type ApiHandler struct {
	engine   *gin.Engine
	dbClient database.DbClient
}

func InstanceHandler() http.Handler {

	ginEngine := gin.Default()
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	mongoDbClient := database.InstanceMongoDb()

	ah := ApiHandler{
		engine:   ginEngine,
		dbClient: mongoDbClient,
	}
	ah.RegisterRoutes()

	return ah
}

func (h ApiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.engine.ServeHTTP(w, r)
}

func (h *ApiHandler) RegisterRoutes() {
	h.engine.GET("/", h.HelloWorldHandler)

	h.engine.GET("/health", h.HealthHandler)

	h.engine.GET("/create_alert/test", h.CreateAlertHandler)
	h.engine.GET("/delete_alerts/test", h.DeleteAlertsHandler)

	h.engine.GET("/create_property", h.CreatePropertyHandler)
	h.engine.GET("/create_properties", h.CreatePropertiesHandler)
	h.engine.GET("/count_properties", h.CountPropertiesHandler)
	h.engine.GET("/delete_properties", h.DeletePropertiesHandler) // To delete testing data

	h.engine.GET("/search", h.SearchHandler) // TODO continue here!!!

	h.engine.Static("/images", "./images")
}
