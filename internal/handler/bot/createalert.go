package telegram

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/venturarome/DaftWatch/internal/model"
	"github.com/venturarome/DaftWatch/internal/scraper"
	"github.com/venturarome/DaftWatch/internal/utils"
)

var searchTypeOptions = []string{"Buy", "Rent"}
var locationOptions = []string{"Dublin", "Cork", "Limerick", "Galway"}
var maxPriceOptions = map[string][]string{
	"Buy":  {"100000", "150000", "200000", "250000", "300000", "350000", "400000", "450000", "500000", "600000", "700000", "800000"},
	"Rent": {"750", "1000", "1250", "1500", "1750", "2000", "2250", "2500", "2750", "3000", "3250", "3500"},
}
var minBedroomsOptions = []string{"1", "2", "3", "4"}

func (th *TelegramHandler) HandleCreateAlert(update tgbotapi.Update) (msg tgbotapi.MessageConfig, clearContext bool) {
	if update.Message == nil {
		panic("Received non-Message Update")
	}

	messageText := strings.Trim(update.Message.Text, " ")
	if strings.Index(messageText, "/createalert") != 0 {
		panic("Badly routed Message Update")
	}

	fmt.Println("[DEBUG] createalert > received: ", messageText)

	userId := update.Message.From.ID
	chatId := update.Message.Chat.ID
	msg = tgbotapi.NewMessage(chatId, "")

	commandParts := strings.Split(messageText, " ")
	switch len(commandParts) {
	case 1:
		// /createalert
		msg.Text = "What are you looking for?"
		msg.ReplyMarkup = CreateKeyboard(searchTypeOptions, 2)
	case 2:
		// /createalert <searchType>
		// TODO validate searchType
		msg.Text = "In which city are you looking for?"
		msg.ReplyMarkup = CreateKeyboard(locationOptions, 2)
	case 3:
		// /createalert <searchType> <location>
		// TODO validate location
		msg.Text = "How much are you willing to spend?"
		msg.ReplyMarkup = CreateKeyboard(maxPriceOptions[commandParts[1]], 2)
	case 4:
		// /createalert <searchType> <location> <maxPrice>
		// TODO validate maxPrice
		msg.Text = "Which is the minimum number of bedrooms?"
		msg.ReplyMarkup = CreateKeyboard(minBedroomsOptions, 2)
	case 5:
		// /createalert <searchType> <location> <maxPrice> <minBedrooms>
		// TODO validate minBedrooms
		// TODO:
		//  0. Extract all pieces of information
		searchType := strings.ToLower(commandParts[1])
		location := strings.ToLower(commandParts[2])
		maxPrice := strings.ToLower(commandParts[3])
		minBedrooms := strings.ToLower(commandParts[4])

		//  1. Scrape Daft with criteria (if possible, asynchronously)
		criteria := scraper.Criteria{
			SearchType: searchType,
			Location:   location,
			Filters: []scraper.Filter{
				{
					Key:   "maxPrice",
					Value: maxPrice,
				},
				{
					Key:   "minBedrooms",
					Value: minBedrooms,
				},
				{
					Key:   "firstPosted",
					Value: "now-20m", // We force to only check super recent listings (last 20 mins), as only want properties from now on.
				},
			},
		}
		scraper.Scrape(criteria) // TODO probar a poner 'go' al inicio.

		//  2. Create alert in DB.
		user := model.User{
			TelegramUserId: userId,
			TelegramChatId: chatId,
		}
		iMaxPrice, _ := utils.StringToInt(maxPrice)
		iMinBedrooms, _ := utils.StringToInt(minBedrooms)
		alert := model.Alert{
			SearchType:  searchType,
			Location:    location,
			MaxPrice:    iMaxPrice,
			MinBedrooms: iMinBedrooms,
		}

		th.dbClient.CreateAlertForUser(alert, user)

		//  3. Reply with elements matching criteria right now.
		msg.Text = "Great! Alert created! I'll send you a message as soon as a new listing appears!"
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		clearContext = true
	default:
		msg.Text = "An error occurred."
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		clearContext = true
	}
	return
}
