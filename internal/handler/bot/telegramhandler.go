package telegram

import (
	"log"
	"math"
	"sort"
	"strings"

	"github.com/venturarome/DaftWatch/internal/database"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramHandler struct {
	dbClient database.DbClient
}

func InstanceHandler() *TelegramHandler {
	return &TelegramHandler{
		dbClient: database.InstanceMongoDb(),
	}
}

const CommandPartsSeparator string = " <|*|> "

func JoinCommandParts(commandParts ...string) string {
	return strings.Join(commandParts, CommandPartsSeparator)
}

func SplitCommandParts(fullCommand string) []string {
	return strings.Split(fullCommand, CommandPartsSeparator)
}

func CreateKeyboard(choices []string, elemsPerRow int, sorted bool) tgbotapi.ReplyKeyboardMarkup {

	if elemsPerRow < 1 {
		log.Panic("elemsPerRow should be a positive integer")
	}

	if sorted {
		sort.Strings(choices)
	}

	numChoices := len(choices)
	buttonRows := make([][]tgbotapi.KeyboardButton, 0, int(math.Ceil(float64(numChoices)/float64(elemsPerRow))))
	buttons := make([]tgbotapi.KeyboardButton, 0, elemsPerRow)

	var button tgbotapi.KeyboardButton
	for i, choice := range choices {
		button = tgbotapi.NewKeyboardButton(choice)
		buttons = append(buttons, button)

		if (i+1)%elemsPerRow == 0 || (i+1) == numChoices {
			buttonRows = append(buttonRows, buttons)
			buttons = make([]tgbotapi.KeyboardButton, 0, elemsPerRow)
		}
	}

	return tgbotapi.NewReplyKeyboard(buttonRows...)
}
