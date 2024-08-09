package telegram

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (th *TelegramHandler) HandleFeedback(update tgbotapi.Update) (msg tgbotapi.MessageConfig, clearContext bool) {
	if update.Message == nil { // panic on non-Message updates
		panic("Received non-Message Update")
	}

	messageText := strings.Trim(update.Message.Text, " ")
	if strings.Index(messageText, "/feedback") != 0 && strings.Index(messageText, "/contact") != 0 {
		panic("Badly routed Message Update")
	}

	fmt.Println("[DEBUG] feedback/contact > received: ", messageText)

	userId := update.Message.From.ID
	userName := update.Message.From.UserName
	chatId := update.Message.Chat.ID
	msg = tgbotapi.NewMessage(chatId, "")

	commandParts := strings.SplitN(messageText, " ", 2)
	switch len(commandParts) {
	case 1:
		// /feedback
		msg.Text = "Please, write your message and it will be forwarded to the bot developer."
	default:
		// /feedback <message>
		feedbackMessage := commandParts[1]

		devChatId, _ := strconv.ParseInt(os.Getenv("TELEGRAM_ID_GOD_USER"), 10, 64)
		textToDev := fmt.Sprintf(
			"User %d (username: @%s) of DaftWatch has sent feedback:\n\n<i>%s</i>\n\n",
			userId, userName, feedbackMessage)
		msgToDev := tgbotapi.NewMessage(devChatId, textToDev)
		msgToDev.ParseMode = "html"

		// Hack to send message to different user:
		tgBotApi, _ := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
		tgBotApi.Send(msgToDev)

		msg.Text = "Feedback message sent to the developer; thank you for reaching out!"
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		clearContext = true
	}
	return
}
