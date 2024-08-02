package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id             primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	TelegramUserId int64              `json:"telegram_user_id" bson:"telegram_user_id"`
	TelegramChatId int64              `json:"telegram_chat_id" bson:"telegram_chat_id"`
}
