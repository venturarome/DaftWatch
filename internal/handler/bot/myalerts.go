package telegram

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/venturarome/DaftWatch/internal/model"
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

	userId := update.Message.From.ID
	chatId := update.Message.Chat.ID
	msg = tgbotapi.NewMessage(chatId, "")

	user := model.User{
		TelegramUserId: userId,
	}
	alerts := th.dbClient.ListAlertsForUser(user)

	if len(alerts) == 0 {
		msg.Text = "No alerts found. You can use /createalert to set new alerts."
		return msg, true
	}

	for _, alert := range alerts {
		msg.Text += " • " + alert.Format() + "\n"
	}
	return msg, true
}
