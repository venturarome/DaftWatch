package daemon

import (
	"fmt"
	"log"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/venturarome/DaftWatch/internal/database"
	"github.com/venturarome/DaftWatch/internal/model"
	"github.com/venturarome/DaftWatch/internal/scraper"
	"github.com/venturarome/DaftWatch/internal/utils"
)

type Daemon struct {
	dbClient       database.DbClient
	botApi         *tgbotapi.BotAPI
	cycleFrequency time.Duration
}

func InstanceDaemon() *Daemon {
	tgBotApi, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		log.Panic(err)
	}
	tgBotApi.Debug = os.Getenv("TELEGRAM_MODE") != "release"
	log.Printf("Authorized on account %s", tgBotApi.Self.UserName)

	return &Daemon{
		dbClient:       database.InstanceMongoDb(),
		botApi:         tgBotApi,
		cycleFrequency: time.Minute * 5,
	}
}

const DateTimeMicro = time.DateTime + ".000000"

func (daemon *Daemon) Run() {

	ticker := time.NewTicker(daemon.cycleFrequency)
	defer ticker.Stop()

	for start := range ticker.C {
		if os.Getenv("DAEMON_MODE") == "debug" {
			fmt.Println("[", start.Format(DateTimeMicro), "] Looping over all Alerts...")
		}

		// Loop over all alerts
		alerts := daemon.dbClient.ListAlerts()
		for _, alert := range alerts {
			// 1. Scrape recent properties
			criteria := scraper.CreateCriteriaFromAlert(alert)
			scrapedProperties := scraper.Scrape(criteria)

			// 2. Compare scraped properties with stored properties
			storedProperties := alert.Properties
			newProperties := utils.DiffSlice(
				scrapedProperties,
				storedProperties,
				func(p1 model.Property, p2 model.Property) bool {
					return p1.ListingId == p2.ListingId
				},
			)
			if len(newProperties) == 0 {
				continue
			}

			// 3. Store new properties
			daemon.dbClient.CreateProperties(newProperties)

			// 4. Update alert properties
			daemon.dbClient.SetPropertiesToAlert(alert, scrapedProperties)

			// 5. Notify alert subscribers
			// 5.1. Prepare text
			msgText := fmt.Sprintf("New listings matched your alert to <b>%s</b>\n", alert.Format())
			for _, property := range newProperties {
				msgText += fmt.Sprintf(
					"\n <u>%s</u> \n • Type: %s\n • Price: %d€\n • Bedrooms: %d\n • See in: %s\n",
					property.Address,
					property.Type,
					property.Price,
					property.NumSingleBedrooms+property.NumDoubleBedrooms,
					property.Url,
				)
			}
			// 5.2. Send Message
			for _, user := range alert.Subscribers {
				msg := tgbotapi.NewMessage(user.TelegramChatId, msgText)
				msg.ParseMode = "html"
				daemon.botApi.Send(msg)
			}
		}

		if os.Getenv("DAEMON_MODE") == "debug" {
			fmt.Println("... Looped over", len(alerts), "Alerts in ", time.Since(start))
		}
	}
}
