package telegram

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// TODO maybe it's a good idea to rename the parent folder to 'process'

//var createAlertSteps = InitSteps([]string{"searchType", "location", "maxPrice"})

// TODO create a util func to create the keyboard given the list of strings and an array with the buttons per row

var searchTypeKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Buy"),
		tgbotapi.NewKeyboardButton("Rent"),
	),
)

var locationKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Dublin"),
		tgbotapi.NewKeyboardButton("Cork"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Limerick"),
		tgbotapi.NewKeyboardButton("Galway"),
	),
)

var maxPriceKeyboard = map[string]tgbotapi.ReplyKeyboardMarkup{
	"Buy": tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("100000"),
			tgbotapi.NewKeyboardButton("150000"),
			tgbotapi.NewKeyboardButton("200000"),
			tgbotapi.NewKeyboardButton("250000"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("300000"),
			tgbotapi.NewKeyboardButton("350000"),
			tgbotapi.NewKeyboardButton("400000"),
			tgbotapi.NewKeyboardButton("450000"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("500000"),
			tgbotapi.NewKeyboardButton("600000"),
			tgbotapi.NewKeyboardButton("700000"),
			tgbotapi.NewKeyboardButton("800000"),
		),
	),
	"Rent": tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("750"),
			tgbotapi.NewKeyboardButton("1000"),
			tgbotapi.NewKeyboardButton("1250"),
			tgbotapi.NewKeyboardButton("1500"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("1750"),
			tgbotapi.NewKeyboardButton("2000"),
			tgbotapi.NewKeyboardButton("2250"),
			tgbotapi.NewKeyboardButton("2500"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("2750"),
			tgbotapi.NewKeyboardButton("3000"),
			tgbotapi.NewKeyboardButton("3250"),
			tgbotapi.NewKeyboardButton("3500"),
		),
	),
}

func HandleCreateAlert(bot *tgbotapi.BotAPI, update tgbotapi.Update) tgbotapi.MessageConfig {
	if update.Message == nil { // panic on non-Message updates
		panic("Received non-Message Update")
	}

	messageText := strings.Trim(update.Message.Text, " ")
	if strings.Index(messageText, "/createalert") != 0 {
		panic("Badly routed Message Update")
	}

	fmt.Println("[DEBUG] createalert > received: ", update.Message.Text)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

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
