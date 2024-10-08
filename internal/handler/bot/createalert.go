package telegram

import (
	"fmt"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/venturarome/DaftWatch/internal/model"
	"github.com/venturarome/DaftWatch/internal/scraper"
	"github.com/venturarome/DaftWatch/internal/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var searchTypeOptions = utils.MapKeysToSlice(scraper.SearchTypesMap) //[]string{"Buy", "Rent"}
var locationOptions = utils.MapKeysToSlice(scraper.LocationsMap)     //[]string{"Belfast", "Cork", "Dublin 1", "Limerick", "Galway"}
var maxPriceOptions = map[string][]string{
	"Buy":  {"100000", "150000", "200000", "250000", "300000", "350000", "400000", "450000", "500000", "600000", "700000", "800000"},
	"Rent": {"750", "1000", "1250", "1500", "1750", "2000", "2250", "2500", "2750", "3000", "3250", "3500"},
}
var minBedroomsOptions = []string{"1", "2", "3", "4"}

const maxAlertsPerUser int = 3

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
	userName := update.Message.From.UserName
	chatId := update.Message.Chat.ID
	msg = tgbotapi.NewMessage(chatId, "")

	commandParts := SplitCommandParts(messageText)
	switch len(commandParts) {
	case 1:
		// /createalert

		// Fair usage
		user := model.User{TelegramUserId: userId}
		if userId != int64(utils.StringToInt(os.Getenv("TELEGRAM_ID_GOD_USER"))) &&
			len(th.dbClient.ListAlertsForUser(user)) >= maxAlertsPerUser {
			msg.Text = fmt.Sprintf("To guarantee a fair use, a limit of %d simultaneous alerts per user is in place. "+
				"However, if you wish to have it extended, consider donating by typing /donate to have it unlocked.", maxAlertsPerUser)
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			clearContext = true
			return
		}

		msg.Text = "What are you looking for?"
		msg.ReplyMarkup = CreateKeyboard(searchTypeOptions, 2, true)
	case 2:
		// /createalert <searchType>
		msg.Text = "Where are you looking for?"
		msg.ReplyMarkup = CreateKeyboard(locationOptions, 2, true)
	case 3:
		// /createalert <searchType> <location>
		msg.Text = "How much are you willing to spend?"
		msg.ReplyMarkup = CreateKeyboard(maxPriceOptions[commandParts[1]], 2, true)
	case 4:
		// /createalert <searchType> <location> <maxPrice>
		msg.Text = "Which is the minimum number of bedrooms?"
		msg.ReplyMarkup = CreateKeyboard(minBedroomsOptions, 2, true)
	case 5:
		// /createalert <searchType> <location> <maxPrice> <minBedrooms>

		// 0. Prepare reply message
		msg.Text = "Great! Alert created! I'll send you a message as soon as a new listing appears!"
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		clearContext = true

		// 1. Create entities
		user := model.User{
			TelegramUserId:   userId,
			TelegramUserName: userName,
			TelegramChatId:   chatId,
		}
		alert := model.Alert{
			SearchType:  commandParts[1],
			Location:    commandParts[2],
			MaxPrice:    utils.StringToInt(commandParts[3]),
			MinBedrooms: utils.StringToInt(commandParts[4]),
		}
		criteria := scraper.CreateCriteriaFromAlert(alert)

		// 2. DB and scraping actions
		go th.doHandleCreateAlert(user, alert, criteria)

	default:
		msg.Text = "An error occurred."
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		clearContext = true
	}
	return
}

func (th *TelegramHandler) doHandleCreateAlert(user model.User, alert model.Alert, criteria scraper.Criteria) {
	// 1. Create User if missing in DB
	th.dbClient.CreateUser(user)

	// 2. Add User as subscriber to Alert in DB
	res := th.dbClient.AddSubscriberToAlert(alert, user) // <-- this upserts the Alert
	if res["UpsertedCount"] == 0 {
		// Nothing else to do if the alert exists already
		return
	} else {
		alert.Id = res["UpsertedID"].(primitive.ObjectID)
	}

	// 3. Scrape Properties matching the Criteria
	scrapedProperties := scraper.Scrape(criteria)

	// 4. Create Properties in DB
	th.dbClient.CreateProperties(scrapedProperties)

	// 5. Add properties' listing IDs to Alert in DB
	th.dbClient.SetPropertiesToAlert(alert, scrapedProperties)
}
