package telegram

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// TODO maybe it's a good idea to rename the parent folder to 'process'

func HandleMyAlerts(bot *tgbotapi.BotAPI, update tgbotapi.Update) tgbotapi.MessageConfig {
	if update.Message == nil { // panic on non-Message updates
		panic("Received non-Message Update")
	}

	messageText := strings.Trim(update.Message.Text, " ")
	if strings.Index(messageText, "/myalerts") != 0 {
		panic("Badly routed Message Update")
	}

	fmt.Println("[DEBUG] myalerts > received: ", update.Message.Text)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	// TODO go to DB and search for all active alerts
	//Alert
	commandParts := strings.Split(update.Message.Text, " ")
	switch len(commandParts) {
	case 1:
		// /createalert
		msg.Text = "What are you looking for?"
		msg.ReplyMarkup = searchTypeKeyboard
	case 2:
		// /createalert <searchType>
		// TODO validate searchType
		msg.Text = "Where are you looking for?"
		msg.ReplyMarkup = locationKeyboard
	case 3:
		// /createalert <searchType> <location>
		// TODO validate location
		msg.Text = "How much are you willing to spend?"
		msg.ReplyMarkup = maxPriceKeyboard[commandParts[1]]
	case 4:
		// /createalert <searchType> <location> <maxPrice>
		// TODO validate location
		// TODO:
		//  1. Scrape Daft with criteria.
		//  2. Create alert in DB.
		//  3. Reply with elements matching criteria right now.
		msg.Text = "Great! Alert created! I'll send you a message as soon as a new listing appears!"
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	}

	return msg
}
