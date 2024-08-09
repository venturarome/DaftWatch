package telegram

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var godCommandsKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("cleanDB"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("cancel"),
	),
)

func (th *TelegramHandler) HandleGodMode(update tgbotapi.Update) (msg tgbotapi.MessageConfig, clearContext bool) {
	if update.Message == nil {
		panic("Received non-Message Update")
	}

	messageText := strings.Trim(update.Message.Text, " ")
	if strings.Index(messageText, "/godmode") != 0 {
		panic("Badly routed Message Update")
	}

	fmt.Println("[DEBUG] godmode > received: ", messageText)

	userId := update.Message.From.ID
	chatId := update.Message.Chat.ID

	godUserId, _ := strconv.ParseInt(os.Getenv("TELEGRAM_ID_GOD_USER"), 10, 64)
	if userId != godUserId {
		msg = tgbotapi.NewMessage(chatId, "You are not God (￢_￢)")
		return
	}
	msg = tgbotapi.NewMessage(chatId, "")

	commandParts := SplitCommandParts(messageText)
	switch len(commandParts) {
	case 1:
		// /godmode
		msg.Text = "Beloved Lord, command me!"
		msg.ReplyMarkup = godCommandsKeyboard
	case 2:
		// /godmode <godCommand>
		// TODO validate godCommand
		godCommand := strings.ToLower(commandParts[1])
		switch godCommand {
		case "cleandb":
			th.dbClient.DeleteAlerts()
			th.dbClient.DeleteProperties()
			th.dbClient.DeleteUsers()

			msg.Text = "Done!"
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			clearContext = true

		case "test":
			// START TEST
			msg.Text = "Blueprint for tests"
			// END TEST

			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			clearContext = true

		case "cancel":
			msg.Text = "I'll be here at anytime, My Lord."
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			clearContext = true

		default:
			msg.Text = "I don't understand your command, My Lord."
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			clearContext = true
		}
	default:
		msg.Text = "An error occurred."
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		clearContext = true
	}
	return
}
