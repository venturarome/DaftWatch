package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id               primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	TelegramChatId   int64              `json:"telegram_chat_id" bson:"telegram_chat_id"`
	TelegramUserId   int64              `json:"telegram_user_id" bson:"telegram_user_id"`
	TelegramUserName string             `json:"telegram_user_name" bson:"telegram_user_name"`
}
