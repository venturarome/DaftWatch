package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id             primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	TelegramId     int64              `json:"telegram_id" bson:"telegram_id"`
	TelegramChatId int64              `json:"telegram_chat_id" bson:"telegram_chat_id"`
}
