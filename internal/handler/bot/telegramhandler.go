package telegram

import (
	"github.com/venturarome/DaftWatch/internal/database"
)

type TelegramHandler struct {
	dbClient database.DbClient
}

func InstanceHandler() *TelegramHandler {
	return &TelegramHandler{
		dbClient: database.InstanceMongoDb(),
	}
}
