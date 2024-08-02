package api

import (
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/venturarome/DaftWatch/internal/model"
)

func (h *ApiHandler) CreateAlertHandler(c *gin.Context) {
	// For testing purposes
	alert := model.Alert{
		SearchType:  "buy",
		Location:    "dublin",
		MaxPrice:    350000,
		MinBedrooms: 3,
	}

	user := model.User{
		TelegramUserId: 123456789,
		TelegramChatId: 987654320 + int64(rand.Int())%5,
	}

	c.JSON(http.StatusOK, h.dbClient.CreateAlertForUser(alert, user))
}
