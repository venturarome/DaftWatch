package telegram

import (
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/venturarome/DaftWatch/internal/model"
	"github.com/venturarome/DaftWatch/internal/utils"
)

var alertToDeleteOptions = func(numAlerts int) []string {
	var options []string
	for i := 1; i <= numAlerts; i++ {
		options = append(options, strconv.Itoa(i))
	}
	return options
}

func (th *TelegramHandler) HandleDeleteAlert(update tgbotapi.Update) (msg tgbotapi.MessageConfig, clearContext bool) {
	if update.Message == nil {
		panic("Received non-Message Update")
	}

	messageText := strings.Trim(update.Message.Text, " ")
	if strings.Index(messageText, "/deletealert") != 0 {
		panic("Badly routed Message Update")
	}

	fmt.Println("[DEBUG] deletealert > received: ", messageText)

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

	commandParts := strings.Split(messageText, " ")
	switch len(commandParts) {
	case 1:
		// /deletealert
		msg.Text = "These are your alerts:\n"
		for i, alert := range alerts {
			msg.Text += strconv.Itoa(i+1) + ") " + alert.Format() + "\n"
		}
		msg.Text += "Which one do you want to delete?"
		msg.ReplyMarkup = CreateKeyboard(alertToDeleteOptions(len(alerts)), 4)
	case 2:
		// /deletealert <alertNum>
		// TODO validate alertNum
		alertNum := utils.StringToInt(commandParts[1])
		user := model.User{
			TelegramUserId: userId,
			TelegramChatId: chatId,
		}

		th.dbClient.RemoveSubscriberFromAlert(alerts[alertNum-1], user)

		msg.Text = "Alert deleted successfully!"
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		clearContext = true
	default:
		msg.Text = "An error occurred."
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		clearContext = true
	}
	return
}
