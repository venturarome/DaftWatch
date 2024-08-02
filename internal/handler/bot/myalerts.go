package telegram

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (th *TelegramHandler) HandleMyAlerts(update tgbotapi.Update) (msg tgbotapi.MessageConfig, clearContext bool) {
	if update.Message == nil { // panic on non-Message updates
		panic("Received non-Message Update")
	}

	messageText := strings.Trim(update.Message.Text, " ")
	if strings.Index(messageText, "/myalerts") != 0 {
		panic("Badly routed Message Update")
	}

	fmt.Println("[DEBUG] myalerts > received: ", messageText)

	msg = tgbotapi.NewMessage(update.Message.Chat.ID, "")

	// TODO go to DB and search for all active alerts

	return msg, false
}
