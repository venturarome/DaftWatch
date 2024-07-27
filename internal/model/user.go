package model

type User struct {
	Id             string `json:"_id" bson:"_id"`
	TelegramId     int64  `json:"telegram_id" bson:"telegram_id"`
	TelegramChatId int64  `json:"telegram_chat_id" bson:"telegram_chat_id"`
}
